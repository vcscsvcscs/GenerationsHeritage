import { KVNamespace } from '@cloudflare/workers-types';
// See https://svelte.dev/docs/kit/types#app.d.ts
// for information about these interfaces
declare global {
	namespace App {
		interface Locals {
			session: Session | null;
		}
		interface Platform {
			env: {
				GH_MEDIA: R2Bucket;
				GH_SESSIONS: KVNamespace;
			};
			cf: CfProperties;
			ctx: ExecutionContext;
		}
	}
}

export {};
