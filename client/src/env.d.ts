/// <reference path="../.astro/types.d.ts" />
interface ImportMetaEnv {
    readonly PUBLIC_API_AUTH: string;
    readonly PUBLIC_API_CORE: string;
    readonly PUBLIC_API_REPORT: string;
    readonly PUBLIC_API_SMART: string;
}

interface ImportMeta {
    readonly env: ImportMetaEnv;
}