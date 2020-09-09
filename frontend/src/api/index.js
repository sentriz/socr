const base_url = process.env.API_BASE;

export const doSearch = async (body) => {
  const response = await fetch(`${base_url}/search`, {
    method: "POST",
    body: JSON.stringify(body),
  });
  return await response.json();
};
