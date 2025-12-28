import { communityThemes, getThemeById } from "./themes";

describe("communityThemes", () => {
  it("exposes at least one theme", () => {
    expect(communityThemes.length).toBeGreaterThan(0);
  });

  it("returns a theme by id", () => {
    const first = communityThemes[0];
    expect(first).toBeDefined();
    const fetched = getThemeById(first.id);
    expect(fetched?.name).toBe(first.name);
    expect(fetched?.previewHtml).toContain("<style>");
  });
});
