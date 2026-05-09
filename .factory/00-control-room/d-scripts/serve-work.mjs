#!/usr/bin/env node
import { spawnSync } from "node:child_process";
import fs from "node:fs";
import os from "node:os";
import path from "node:path";
import process from "node:process";

const scriptName = "serve-work.mjs";

function candidates(repoRoot) {
  const roots = [];
  if (process.env.FACTORY_SKILL_ROOT) roots.push(process.env.FACTORY_SKILL_ROOT);
  roots.push(
    path.join(repoRoot, ".agents/skills/factory"),
    path.join(repoRoot, ".codex/skills/factory"),
    path.join(process.env.CODEX_HOME || path.join(os.homedir(), ".codex"), "skills/factory"),
    path.join(os.homedir(), ".agents/skills/factory")
  );
  return roots.map((root) => path.join(root, "scripts", scriptName));
}

const repoRootResult = spawnSync("git", ["rev-parse", "--show-toplevel"], { encoding: "utf8" });
if (repoRootResult.status !== 0) {
  console.error("Unable to resolve repository root.");
  process.exit(1);
}

for (const candidate of candidates(repoRootResult.stdout.trim())) {
  if (fs.existsSync(candidate)) {
    const result = spawnSync(process.execPath, [candidate, ...process.argv.slice(2)], {
      stdio: "inherit"
    });
    process.exit(result.status || 0);
  }
}

console.error(`Unable to locate executable factory skill script: scripts/${scriptName}`);
process.exit(1);
