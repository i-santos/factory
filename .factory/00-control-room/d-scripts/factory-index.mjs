#!/usr/bin/env node
import crypto from "node:crypto";
import fs from "node:fs";
import path from "node:path";

const args = process.argv.slice(2);

function readArg(name, fallback) {
  const index = args.indexOf(`--${name}`);
  if (index >= 0 && args[index + 1]) return args[index + 1];
  return fallback;
}

const mode = readArg("mode", args[0] || "status");
const workspaceRoot = readArg("workspace-root", ".factory");
const indexDir = readArg(
  "index-dir",
  path.join(workspaceRoot, "00-control-room", "g-index")
);

if (!["rebuild", "verify", "status"].includes(mode)) {
  console.error(JSON.stringify({
    status: "blocked",
    issue_class: "invalid-index-mode",
    allowed_modes: ["rebuild", "verify", "status"]
  }));
  process.exit(2);
}

const sourceRoots = [
  "00-control-room/a-config",
  "01-dock",
  "02-yard",
  "03-shop-floor",
  "04-finished-goods",
  "05-labs"
];

function toPosix(value) {
  return value.split(path.sep).join("/");
}

function readText(filePath) {
  return fs.readFileSync(filePath, "utf8");
}

function stableJson(value) {
  return JSON.stringify(sortValue(value));
}

function sortValue(value) {
  if (Array.isArray(value)) return value.map(sortValue);
  if (value && typeof value === "object") {
    return Object.keys(value).sort().reduce((result, key) => {
      result[key] = sortValue(value[key]);
      return result;
    }, {});
  }
  return value;
}

function hashText(value) {
  return crypto.createHash("sha256").update(value).digest("hex");
}

function listFiles(root) {
  if (!fs.existsSync(root)) return [];
  const stat = fs.statSync(root);
  if (stat.isFile()) return [root];
  const entries = fs.readdirSync(root, { withFileTypes: true }).sort((a, b) => a.name.localeCompare(b.name));
  const files = [];
  for (const entry of entries) {
    const fullPath = path.join(root, entry.name);
    if (entry.isDirectory()) {
      files.push(...listFiles(fullPath));
    } else if (entry.isFile()) {
      files.push(fullPath);
    }
  }
  return files;
}

function canonicalFiles() {
  return sourceRoots
    .flatMap((root) => listFiles(path.join(workspaceRoot, root)))
    .filter((filePath) => !toPosix(filePath).includes("/00-control-room/g-index/"))
    .filter((filePath) => !toPosix(filePath).includes("/00-control-room/e-state/"))
    .filter((filePath) => !toPosix(filePath).includes("/00-control-room/f-logs/"))
    .filter((filePath) => !filePath.endsWith("/.gitkeep"))
    .sort((a, b) => toPosix(a).localeCompare(toPosix(b)));
}

function fingerprint(files) {
  const entries = files.map((filePath) => {
    const relativePath = toPosix(path.relative(workspaceRoot, filePath));
    return {
      path: relativePath,
      sha256: hashText(readText(filePath))
    };
  });
  return hashText(stableJson(entries));
}

function firstMatch(text, patterns) {
  for (const pattern of patterns) {
    const match = text.match(pattern);
    if (match) return match[1].trim();
  }
  return null;
}

function inferTimestamp(recordId) {
  const match = recordId.match(/^(\d{8}T\d{6}Z)/);
  return match ? match[1] : null;
}

function inferLane(relativePath) {
  const parts = relativePath.split("/");
  if (parts[0] === "00-control-room") return parts.slice(0, 3).join("/");
  return parts.slice(0, 2).join("/");
}

function inferRecordType(relativePath) {
  if (relativePath.startsWith("00-control-room/a-config/")) return "config";
  if (relativePath.startsWith("01-dock/e-deferred/")) return "deferred-item";
  if (relativePath.startsWith("01-dock/")) return "dock-record";
  if (relativePath.startsWith("02-yard/")) return "work-package";
  if (relativePath.startsWith("03-shop-floor/a-input-buffer/")) return "work-order";
  if (relativePath.startsWith("03-shop-floor/")) return "shop-floor-record";
  if (relativePath.startsWith("04-finished-goods/")) return "finished-good";
  if (relativePath.startsWith("05-labs/")) return "lab-artifact";
  return "factory-record";
}

function inferReleaseState(relativePath, text) {
  if (relativePath.startsWith("04-finished-goods/b-release-ready/")) return "release-ready";
  if (relativePath.startsWith("04-finished-goods/c-shipped/")) return "shipped";
  if (relativePath.startsWith("04-finished-goods/c-released/")) return "released-environment";
  if (relativePath.startsWith("04-finished-goods/a-awaiting-release-qc/")) return "awaiting-release-qc";
  return firstMatch(text, [
    /^release_state:\s*"?([^"\n]+)"?/mi,
    /^Release state:\s*`?([^`\n]+)`?/mi
  ]);
}

function inferStatus(relativePath, text) {
  return firstMatch(text, [
    /^status:\s*"?([^"\n]+)"?/mi,
    /^Status:\s*`?([^`\n]+)`?/mi,
    /^- Status:\s*`?([^`\n]+)`?/mi
  ]) || inferReleaseState(relativePath, text);
}

function extractRelatedPaths(text) {
  const paths = new Set();
  const pattern = /\.factory\/[A-Za-z0-9._/\-]+/g;
  for (const match of text.matchAll(pattern)) {
    paths.add(match[0].replace(/[),.;]+$/, ""));
  }
  return [...paths].sort();
}

function extractVerificationEvidence(text) {
  const lines = text.split(/\r?\n/);
  const evidence = [];
  let inVerification = false;
  for (const line of lines) {
    if (/^#{1,4}\s+Verification\b/i.test(line)) {
      inVerification = true;
      continue;
    }
    if (inVerification && /^#{1,4}\s+/.test(line)) break;
    if (inVerification && line.trim()) evidence.push(line.trim());
  }
  return evidence.slice(0, 20);
}

function buildRecord(filePath) {
  const relativePath = toPosix(path.relative(workspaceRoot, filePath));
  const text = readText(filePath);
  const basename = path.basename(filePath);
  const recordId = basename.replace(/\.(md|json|jsonl|txt)$/i, "");
  const title = firstMatch(text, [/^#\s+(.+)$/m, /^title:\s*"?([^"\n]+)"?/mi]);
  const packageBranch = firstMatch(text, [
    /^package_branch:\s*"?([^"\n]+)"?/mi,
    /^Package branch:\s*`([^`]+)`/mi,
    /^- Package branch:\s*`([^`]+)`/mi,
    /^Shipping branch:\s*`([^`]+)`/mi
  ]);
  const baseBranch = firstMatch(text, [
    /^base_branch:\s*"?([^"\n]+)"?/mi,
    /^Base branch:\s*`([^`]+)`/mi,
    /^- Base branch:\s*`([^`]+)`/mi
  ]);
  const releaseState = inferReleaseState(relativePath, text);
  const status = inferStatus(relativePath, text);
  const warnings = [];
  if (relativePath.startsWith("04-finished-goods/b-release-ready/") && status && !String(status).includes("release-ready")) {
    warnings.push("release-ready-lane-status-mismatch");
  }
  if (relativePath.startsWith("04-finished-goods/c-shipped/") && !status) {
    warnings.push("shipped-record-missing-status");
  }
  if (relativePath.startsWith("04-finished-goods/c-released/") && relativePath.split("/").length < 4) {
    warnings.push("released-record-missing-environment-folder");
  }
  return {
    derived: true,
    id: recordId,
    type: inferRecordType(relativePath),
    lane: inferLane(relativePath),
    path: `${workspaceRoot}/${relativePath}`,
    title,
    outcome: firstMatch(text, [/^##\s+Outcome\s*\n+([^\n]+)/mi, /^outcome:\s*"?([^"\n]+)"?/mi]),
    lifecycle_status: status,
    package_branch: packageBranch,
    base_branch: baseBranch,
    source: firstMatch(text, [/^source:\s*"?([^"\n]+)"?/mi, /^Source:\s*`?([^`\n]+)`?/mi]),
    provenance: firstMatch(text, [/^provenance:\s*"?([^"\n]+)"?/mi, /^Provenance:\s*`?([^`\n]+)`?/mi]),
    related_paths: extractRelatedPaths(text),
    release_state: releaseState,
    verification_evidence: extractVerificationEvidence(text),
    timestamp: inferTimestamp(recordId),
    warnings
  };
}

function buildIndex() {
  const files = canonicalFiles();
  const sourceFingerprint = fingerprint(files);
  const generatedAt = new Date().toISOString();
  const records = files.map(buildRecord);
  const warnings = records.flatMap((record) => record.warnings.map((warning) => ({
    record_id: record.id,
    path: record.path,
    warning
  })));
  const countsByLane = records.reduce((counts, record) => {
    counts[record.lane] = (counts[record.lane] || 0) + 1;
    return counts;
  }, {});
  const countsByType = records.reduce((counts, record) => {
    counts[record.type] = (counts[record.type] || 0) + 1;
    return counts;
  }, {});
  const relationships = records
    .filter((record) => record.related_paths.length > 0 || record.source || record.provenance)
    .map((record) => ({
      derived: true,
      id: record.id,
      path: record.path,
      source: record.source,
      provenance: record.provenance,
      related_paths: record.related_paths
    }));
  const releaseLedger = records
    .filter((record) => record.path.includes("/04-finished-goods/"))
    .map((record) => ({
      derived: true,
      id: record.id,
      path: record.path,
      title: record.title,
      release_state: record.release_state,
      lifecycle_status: record.lifecycle_status,
      package_branch: record.package_branch,
      base_branch: record.base_branch,
      timestamp: record.timestamp,
      warnings: record.warnings
    }));
  const timeline = records
    .filter((record) => record.timestamp)
    .sort((a, b) => a.timestamp.localeCompare(b.timestamp))
    .map((record) => ({
      derived: true,
      timestamp: record.timestamp,
      id: record.id,
      type: record.type,
      lane: record.lane,
      path: record.path,
      title: record.title,
      lifecycle_status: record.lifecycle_status,
      release_state: record.release_state
    }));
  return {
    manifest: {
      derived: true,
      generated_at: generatedAt,
      source: "canonical .factory records",
      source_roots: sourceRoots.map((root) => `${workspaceRoot}/${root}`),
      source_fingerprint: sourceFingerprint,
      record_count: records.length,
      warning_count: warnings.length,
      canonical_files_win_on_conflict: true
    },
    records,
    currentState: {
      derived: true,
      generated_at: generatedAt,
      source_fingerprint: sourceFingerprint,
      record_count: records.length,
      counts_by_lane: countsByLane,
      counts_by_type: countsByType,
      warnings
    },
    relationships: {
      derived: true,
      generated_at: generatedAt,
      source_fingerprint: sourceFingerprint,
      relationships
    },
    timeline,
    releaseLedger
  };
}

function writeJson(filePath, value) {
  fs.writeFileSync(filePath, `${JSON.stringify(value, null, 2)}\n`);
}

function writeJsonl(filePath, rows) {
  fs.writeFileSync(filePath, rows.map((row) => JSON.stringify(row)).join("\n") + (rows.length ? "\n" : ""));
}

function readManifest() {
  const manifestPath = path.join(indexDir, "manifest.json");
  if (!fs.existsSync(manifestPath)) return null;
  return JSON.parse(readText(manifestPath));
}

function inspectFreshness() {
  const files = canonicalFiles();
  const currentFingerprint = fingerprint(files);
  const manifest = readManifest();
  if (!manifest) {
    return {
      status: "missing",
      derived: true,
      index_dir: indexDir,
      current_source_fingerprint: currentFingerprint,
      recovery: "Run factory index with mode rebuild."
    };
  }
  const fresh = manifest.source_fingerprint === currentFingerprint;
  return {
    status: fresh ? "fresh" : "stale",
    derived: true,
    index_dir: indexDir,
    generated_at: manifest.generated_at,
    indexed_source_fingerprint: manifest.source_fingerprint,
    current_source_fingerprint: currentFingerprint,
    record_count: manifest.record_count,
    warning_count: manifest.warning_count,
    canonical_files_win_on_conflict: true,
    recovery: fresh ? null : "Run factory index with mode rebuild."
  };
}

if (mode === "rebuild") {
  const index = buildIndex();
  fs.mkdirSync(indexDir, { recursive: true });
  writeJson(path.join(indexDir, "manifest.json"), index.manifest);
  writeJsonl(path.join(indexDir, "records.jsonl"), index.records);
  writeJson(path.join(indexDir, "current-state.json"), index.currentState);
  writeJson(path.join(indexDir, "relationships.json"), index.relationships);
  writeJsonl(path.join(indexDir, "timeline.jsonl"), index.timeline);
  writeJsonl(path.join(indexDir, "release-ledger.jsonl"), index.releaseLedger);
  console.log(JSON.stringify({
    status: "rebuilt",
    derived: true,
    index_dir: indexDir,
    source_fingerprint: index.manifest.source_fingerprint,
    record_count: index.manifest.record_count,
    warning_count: index.manifest.warning_count,
    files: [
      "manifest.json",
      "records.jsonl",
      "current-state.json",
      "relationships.json",
      "timeline.jsonl",
      "release-ledger.jsonl"
    ].map((file) => path.join(indexDir, file)),
    canonical_files_win_on_conflict: true
  }, null, 2));
} else {
  const freshness = inspectFreshness();
  console.log(JSON.stringify({
    mode,
    ...freshness
  }, null, 2));
}
