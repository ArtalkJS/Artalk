type Graph = Map<string, Set<string>> // Graph with nodes and their dependencies
type SCC = Set<string>[] // Array of Strongly Connected Components

/**
 * Tarjan's Algorithm to find Strongly Connected Components (SCCs) in a directed graph.
 *
 * The function uses a depth-first search (DFS) approach to discover SCCs in the input graph.
 * It keeps track of discovery times and low-link values for each node, which helps in detecting cycles.
 * Once SCCs are found, they are added as sets to the result array.
 *
 * @param graph - A directed graph represented as a Map where the keys are nodes and values are arrays of adjacent nodes.
 * @returns An array of Sets, where each Set contains the nodes of one Strongly Connected Component.
 * @link https://en.wikipedia.org/wiki/Tarjan%27s_strongly_connected_components_algorithm#The_algorithm_in_pseudocode
 */
export function tarjan(graph: Graph): SCC {
  const indices = new Map<string, number>() // To store the discovery index of each node
  const lowlinks = new Map<string, number>() // To store the lowest point reachable from each node
  const onStack = new Set<string>() // To keep track of nodes currently on the stack
  const stack: string[] = [] // Stack to simulate recursion and track SCC nodes
  const scc: SCC = [] // Result array to store SCCs
  let idx = 0 // Global index counter

  /**
   * Strongly connects a node by performing DFS, updating the indices and lowlinks.
   * Once an SCC is found (when a node's lowlink equals its index), it's popped off the stack.
   *
   * @param v - The current node being explored in the DFS
   */
  function strongConnect(v: string): void {
    // Set the discovery index and lowlink for the node
    indices.set(v, idx)
    lowlinks.set(v, idx)
    idx++
    stack.push(v)
    onStack.add(v)

    // Explore the neighbors (dependencies) of the current node
    const deps = graph.get(v) || [] // Get the adjacent nodes (or an empty array if no edges)
    for (const dep of deps) {
      if (!indices.has(dep)) {
        // If the neighbor hasn't been visited, recursively explore it
        strongConnect(dep)
        lowlinks.set(v, Math.min(lowlinks.get(v)!, lowlinks.get(dep)!))
      } else if (onStack.has(dep)) {
        // If the neighbor is on the stack, update the lowlink of the current node
        lowlinks.set(v, Math.min(lowlinks.get(v)!, indices.get(dep)!))
      }
    }

    // If the current node is a root node (its lowlink equals its index), it forms an SCC
    if (lowlinks.get(v) === indices.get(v)) {
      const vertices = new Set<string>()
      let w: string | undefined = undefined
      // Pop all nodes off the stack until we return to the current node
      while (v !== w) {
        w = stack.pop()!
        onStack.delete(w)
        vertices.add(w)
      }
      // Add the SCC to the result
      scc.push(vertices)
    }
  }

  // Start DFS on all nodes that haven't been visited yet
  for (const v of graph.keys()) {
    if (!indices.has(v)) {
      strongConnect(v)
    }
  }

  // Return the list of SCCs
  return scc
}
