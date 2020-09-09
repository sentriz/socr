const baseURL = process.env.SOCR_API_URL;

console.log(baseURL)

export const doSearch = async (body) => {
  const response = await fetch(`${base_url}/search`, {
    method: "POST",
    body: JSON.stringify(body),
  });
  return await response.json();
};
