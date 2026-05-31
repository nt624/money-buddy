/* eslint-disable @typescript-eslint/no-require-imports */
const { createDefaultPreset } = require("ts-jest");

const tsJestTransformCfg = createDefaultPreset().transform;

/** @type {import("jest").Config} **/
module.exports = {
  testEnvironment: "node",
  transform: {
    ...tsJestTransformCfg,
  },
  moduleNameMapper: {
    "^@/lib/firebase/config$": "<rootDir>/src/lib/firebase/__mocks__/config.ts",
    "^@pace/core/types/(.*)$": "<rootDir>/../../packages/core/src/types/$1",
    "^@pace/core/types$": "<rootDir>/../../packages/core/src/types/index.ts",
    "^@pace/core/api$": "<rootDir>/../../packages/core/src/api/index.ts",
    "^@pace/core/auth$": "<rootDir>/../../packages/core/src/auth/index.ts",
    "^@pace/core/data$": "<rootDir>/../../packages/core/src/data/index.ts",
    "^@pace/core/hooks$": "<rootDir>/../../packages/core/src/hooks/index.ts",
    "^@pace/core$": "<rootDir>/../../packages/core/src/index.ts",
    "^@/(.*)$": "<rootDir>/src/$1",
  },
};