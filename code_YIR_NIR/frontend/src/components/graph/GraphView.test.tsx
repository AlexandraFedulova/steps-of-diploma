import { render, screen, waitFor } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import { beforeEach, describe, expect, it, vi } from "vitest";
import { sampleGraphPayload } from "@/test/fixtures/graphPayload";
import {
  createMockCy,
  cytoscapeConstructorMock,
  type MockCy,
} from "@/test/mocks/cytoscapeMock";
import { GraphView } from "./GraphView";

vi.mock("cytoscape", () => ({
  default: cytoscapeConstructorMock,
}));

vi.mock("cytoscape-fcose", () => ({
  default: vi.fn(),
}));

vi.mock("./cytoscapeSetup", () => ({
  registerCytoscapeExtensions: vi.fn(),
}));

vi.mock("@/utils/cytoscapeLayout", () => ({
  buildForceDirectedLayout: vi.fn(() => ({ name: "fcose" })),
}));

vi.mock("@/utils/communityVisual", async (importOriginal) => {
  const actual = await importOriginal<typeof import("@/utils/communityVisual")>();
  return {
    ...actual,
    applyCommunitySeedPositions: vi.fn(),
  };
});

function buildMockNodes() {
  return [
    {
      id: "A",
      data: {
        id: "A",
        name: "Алиса",
        displayLabel: "1",
        communityId: 0,
        communityFill: "#2563eb",
        centralityNorm: 0,
      },
      classes: new Set<string>(),
      selected: false,
      position: { x: 0, y: 0 },
    },
    {
      id: "B",
      data: {
        id: "B",
        name: "Борис",
        displayLabel: "2",
        communityId: 0,
        communityFill: "#2563eb",
        centralityNorm: 0.5,
      },
      classes: new Set<string>(),
      selected: false,
      position: { x: 50, y: 0 },
    },
    {
      id: "C",
      data: {
        id: "C",
        name: "Вера",
        displayLabel: "3",
        communityId: 1,
        communityFill: "#16a34a",
        centralityNorm: 0.2,
      },
      classes: new Set<string>(),
      selected: false,
      position: { x: 100, y: 0 },
    },
  ];
}

describe("GraphView", () => {
  let mockCy: MockCy;

  beforeEach(() => {
    vi.clearAllMocks();
    mockCy = createMockCy(buildMockNodes());
    cytoscapeConstructorMock.mockReturnValue(mockCy as never);
  });

  it("renders graph stats and cytoscape container", async () => {
    render(<GraphView graph={sampleGraphPayload} />);

    expect(screen.getByText(/Узлов:/).closest("span")).toHaveTextContent("3");
    expect(screen.getByText(/Рёбер:/).closest("span")).toHaveTextContent("2");
    expect(screen.getByLabelText("Интерактивная социограмма")).toBeInTheDocument();

    await waitFor(() => {
      expect(cytoscapeConstructorMock).toHaveBeenCalled();
    });

    expect(mockCy.layout).toHaveBeenCalled();
    expect(mockCy.fit).toHaveBeenCalled();
  });

  it("shows tooltip on node hover and hides on mouseout", async () => {
    render(<GraphView graph={sampleGraphPayload} />);

    await waitFor(() => expect(cytoscapeConstructorMock).toHaveBeenCalled());

    mockCy.trigger("mouseover", "B");

    await waitFor(() => {
      expect(screen.getByRole("tooltip")).toBeInTheDocument();
    });
    expect(screen.getByRole("tooltip")).toHaveTextContent(/№2/);
    expect(screen.getByRole("tooltip")).toHaveTextContent(/сообщество/i);

    mockCy.trigger("mouseout", "B");

    await waitFor(() => {
      expect(screen.queryByRole("tooltip")).not.toBeInTheDocument();
    });
  });

  it("shows respondent name in panel on node click", async () => {
    const user = userEvent.setup();
    render(<GraphView graph={sampleGraphPayload} />);

    await waitFor(() => expect(cytoscapeConstructorMock).toHaveBeenCalled());

    mockCy.trigger("tap", "A");

    await waitFor(() => {
      expect(screen.getByRole("complementary", { name: "Информация о респонденте" })).toBeInTheDocument();
    });
    expect(screen.getByRole("heading", { level: 3, name: "Алиса" })).toBeInTheDocument();
    expect(screen.getByText("№ на графе")).toBeInTheDocument();

    await user.click(screen.getByRole("button", { name: "Закрыть" }));

    await waitFor(() => {
      expect(
        screen.queryByRole("complementary", { name: "Информация о респонденте" }),
      ).not.toBeInTheDocument();
    });
  });

  it("destroys cytoscape instance on unmount", async () => {
    const { unmount } = render(<GraphView graph={sampleGraphPayload} />);

    await waitFor(() => expect(cytoscapeConstructorMock).toHaveBeenCalled());
    unmount();

    expect(mockCy.destroy).toHaveBeenCalled();
  });
});
