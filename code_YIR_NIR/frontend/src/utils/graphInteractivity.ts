import type { Core, EventObject, NodeSingular } from "cytoscape";

export interface GraphInteractionHandlers {
  onNodeHover: (nodeId: string | null, position?: { x: number; y: number }) => void;
  onNodeSelect: (nodeId: string | null) => void;
}

export function configureGraphInteraction(cy: Core): void {
  cy.userPanningEnabled(true);
  cy.userZoomingEnabled(true);
  cy.boxSelectionEnabled(false);

  cy.autoungrabify(false);
  cy.nodes().grabify();
  cy.autounselectify(false);
}

export function bindGraphInteractivity(
  cy: Core,
  handlers: GraphInteractionHandlers,
): () => void {
  const clearHoverClass = () => {
    cy.nodes(".hovered").removeClass("hovered");
  };

  const onNodeOver = (event: EventObject) => {
    const node = event.target as NodeSingular;
    node.addClass("hovered");
    const pos = node.renderedPosition();
    handlers.onNodeHover(node.id(), { x: pos.x, y: pos.y });
  };

  const onNodeOut = () => {
    clearHoverClass();
    handlers.onNodeHover(null);
  };

  const onNodeTap = (event: EventObject) => {
    const node = event.target as NodeSingular;
    cy.nodes().unselect();
    node.select();
    handlers.onNodeSelect(node.id());
  };

  const onTap = (event: EventObject) => {
    if (event.target !== cy) {
      return;
    }
    cy.nodes().unselect();
    handlers.onNodeSelect(null);
  };

  const onDragStart = (event: EventObject) => {
    const node = event.target as NodeSingular;
    node.addClass("hovered");
    handlers.onNodeSelect(node.id());
  };

  cy.on("mouseover", "node", onNodeOver);
  cy.on("mouseout", "node", onNodeOut);
  cy.on("tap", "node", onNodeTap);
  cy.on("tap", onTap);
  cy.on("grab", "node", onDragStart);

  return () => {
    cy.off("mouseover", "node", onNodeOver);
    cy.off("mouseout", "node", onNodeOut);
    cy.off("tap", "node", onNodeTap);
    cy.off("tap", onTap);
    cy.off("grab", "node", onDragStart);
    clearHoverClass();
  };
}
