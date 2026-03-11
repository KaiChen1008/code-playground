export interface Language {
  name: string;
  version: string;
}

export interface RunResponse {
  output: string;
  id?: string;
  error?: string;
}

export interface SnippetResponse {
  id: string;
  language: string;
  code: string;
  output: string;
}

export async function fetchLanguages(): Promise<Language[]> {
  const resp = await fetch('/api/v1/languages');
  if (!resp.ok) throw new Error('Failed to fetch languages');
  return resp.json();
}

export async function runCode(language: string, code: string, id?: string): Promise<RunResponse> {
  const body: any = { language, code };
  if (id) body.id = id;
  const resp = await fetch('/api/v1/run', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(body),
  });
  if (!resp.ok) throw new Error('Failed to run code');
  return resp.json();
}

export async function formatCode(language: string, code: string): Promise<{ code: string }> {
  const resp = await fetch('/api/v1/format', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ language, code }),
  });
  if (!resp.ok) throw new Error('Failed to format code');
  return resp.json();
}

export async function loadSnippet(id: string): Promise<SnippetResponse> {
  const resp = await fetch(`/api/v1/snippet/${id}`);
  if (!resp.ok) throw new Error('Failed to load snippet');
  return resp.json();
}

export async function deleteSnippet(id: string): Promise<void> {
  const resp = await fetch(`/api/v1/snippet/${id}`, { method: 'DELETE' });
  if (!resp.ok) throw new Error('Failed to delete snippet');
}
