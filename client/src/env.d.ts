/// <reference path="../.astro/types.d.ts" />
interface ImportMetaEnv {
    readonly PUBLIC_API_AUTH: string;
    readonly PUBLIC_API_TRAFFIC: string;
    readonly PUBLIC_API_REPORT: string;
}

interface ImportMeta {
    readonly env: ImportMetaEnv;
}