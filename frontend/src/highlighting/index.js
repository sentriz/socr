import { fields } from '../api'

export const zipBlocks = (screenshot) => {
  let flatText = screenshot.fields[fields.BLOCKS_TEXT];
  let flatPosition = screenshot.fields[fields.BLOCKS_POSITION];
  if (!Array.isArray(flatText)) flatText = [flatText];
  if (!Array.isArray(flatPosition)) flatPosition = [flatPosition];

  const queriesMatches = screenshot.locations[fields.BLOCKS_TEXT];
  const queryMatches = Object.values(queriesMatches)[0];
  const matchIndexes = new Set(
    queryMatches.map((match) => match.array_positions).flat()
  );

  return flatText.map((block, i) => {
    const [minX, minY, maxX, maxY] = flatPosition.slice(4 * i, 4 * i + 4);
    return {
      text: block,
      position: { minX, minY, maxX, maxY },
      match: matchIndexes.has(i),
    };
  });
};
