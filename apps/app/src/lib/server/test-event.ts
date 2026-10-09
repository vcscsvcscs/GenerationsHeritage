type EventOptions = {
	userId?: number | null;
	params?: Record<string, string>;
	url?: string;
	body?: unknown;
	rawBody?: string;
};

export function makeEvent<E>({
	userId = 7,
	params = {},
	url = 'http://localhost/',
	body,
	rawBody
}: EventOptions = {}): E {
	const payload = rawBody ?? (body === undefined ? undefined : JSON.stringify(body));
	return {
		locals: { session: userId === null ? null : { id: 's', expiresAt: 0, userId } },
		params,
		url: new URL(url),
		request: new Request(url, { method: 'POST', body: payload })
	} as unknown as E;
}

export const upstreamOk = (data?: unknown, status = 200) => ({
	data,
	response: new Response(null, { status })
});

export const upstreamFail = (error: unknown, status: number) => ({
	error,
	response: new Response(null, { status })
});
