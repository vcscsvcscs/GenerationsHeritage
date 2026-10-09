import type { components } from '$lib/api/api.gen';

export async function savePerson(
	personId: string | number | undefined,
	draft: components['schemas']['PersonProperties']
): Promise<string | null> {
	const response = await fetch(`/api/person/${personId}`, {
		method: 'PATCH',
		headers: { 'Content-Type': 'application/json' },
		body: JSON.stringify(draft)
	});
	if (response.ok) return null;

	const details = await response.json().catch(() => null);
	return `Error saving person data, status: ${response.status} ${JSON.stringify(details)}`;
}
