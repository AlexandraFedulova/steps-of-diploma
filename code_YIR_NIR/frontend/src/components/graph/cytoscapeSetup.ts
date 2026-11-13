import cytoscape from "cytoscape";
import fcose from "cytoscape-fcose";

let extensionsRegistered = false;

/** Регистрирует расширения Cytoscape (один раз). */
export function registerCytoscapeExtensions(): void {
  if (extensionsRegistered) {
    return;
  }
  cytoscape.use(fcose);
  extensionsRegistered = true;
}
