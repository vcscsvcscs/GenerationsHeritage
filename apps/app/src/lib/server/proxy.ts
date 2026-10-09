export type ProxyResult = { data?: unknown; error?: unknown; response: Response };

export const jsonResponse = (body: unknown, status: number): Response =>
	new Response(body === undefined ? null : JSON.stringify(body), {
		status,
		headers: { 'Content-Type': 'application/json' }
	});

export const toResponse = ({ data, error, response }: ProxyResult): Response =>
	jsonResponse(response.ok ? data : error, response.status);

export const badRequest = (msg: string): Response => jsonResponse({ msg }, 400);

export async function readJson<T>(request: Request): Promise<T | null> {
	try {
		return (await request.json()) as T;
	} catch {
		return null;
	}
}

export function withSession<E extends { locals: App.Locals }>(
	handler: (event: E, userId: number) => Promise<Response>
): (event: E) => Promise<Response> {
	return (event) =>
		event.locals.session === null
			? Promise.resolve(jsonResponse({ msg: 'Unauthorized' }, 401))
			: handler(event, event.locals.session.userId);
}

export function withSessionAndId<E extends { locals: App.Locals; params: { ID: string } }>(
	handler: (event: E, userId: number, id: number) => Promise<Response>
): (event: E) => Promise<Response> {
	return withSession<E>((event, userId) => {
		const id = Number(event.params.ID);
		return Number.isInteger(id) && id > 0
			? handler(event, userId, id)
			: Promise.resolve(badRequest('Invalid ID'));
	});
}

export const userHeader = (userId: number) => ({ 'X-User-ID': userId });
