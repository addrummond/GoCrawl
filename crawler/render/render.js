const COL_WIDTH = 5;

function renderGraph() {
  let graph = makeInitialGraph();

  const nodeToRow = getNodeRows(GRAPH, NODE_METADATA);

  let i = 0;
  for (const k of Object.keys(GRAPH)) {
    const sk = displayUrl(STRIP_PREFIX, k);
    
    graph.elements.push({
      data: {
        id: k,
        label: sk,
        col: getNodeCol(NODE_METADATA[k].Depth, i),
        row: nodeToRow[k],
        color: nodeColor(NODE_METADATA[k]),
        labelColor: labelColor(NODE_METADATA[k])
      }
    });

    for (const links of GRAPH[k]) {
      graph.elements.push({
        data: {
          source: k,
          target: links.ToUrl,
          color: arrowColor(links.IsAsset)
        }
      });
    }

    ++i;
  }

  cytoscape(graph);
}

function makeInitialGraph() {
  return {
    container: document.body,
    elements: [],
    layout: {
      name: 'grid',
      directed: true,
      position: (node) => {
        return {
          row: node.data('row'),
          col: node.data('col')
        };
      }
    },
    style: [
      {
        selector: 'node',
        style: {
          label: 'data(label)',
          color: 'data(labelColor)',
          'font-size': '50pt',
          'background-color': 'data(color)'
        }
      },
      {
        selector: 'edge',
        style: {
          width: 2,
          'target-arrow-shape': 'triangle',
          'line-color': 'data(color)',
          'target-arrow-color': 'data(color)',
          'curve-style': 'bezier',
          'control-point-step-size': 400,
          'arrow-scale': 3,
        }
      }
    ]
  };
}

function nodeColor(metadata) {
  if (metadata.Depth == 0)
    return 'red';
  if (metadata.PureAsset)
    return 'grey';
  return 'blue';
}

function labelColor(metadata) {
  if (metadata.Depth == 0)
    return 'red';
  if (metadata.PureAsset)
    return 'darkgrey';
  return 'darkblue';
}

function arrowColor(isAsset) {
  if (isAsset)
    return 'grey';
  return 'blue';
}

function displayUrl(prefix, url) {
  if (url == prefix)
    return url;
  if (url.indexOf(prefix) == 0) {
    const rest = url.substr(prefix.length);
    if (rest.length > 0 && rest.charAt(0) != '/')
      return "/" + rest;
    else
      return rest;
  }
  return url;
}

function getNodeRows(graph, nodeMetadata) {
  let withMetadata = Object.keys(graph);
  for (let i = 0; i < withMetadata.length; ++i) {
    withMetadata[i] = [withMetadata[i], nodeMetadata[withMetadata[i]]];
  }
  withMetadata = withMetadata.sort(([url1, md1], [url2, md2]) => {
    if (graph[url1].length == graph[url2].length)
      return md2.Popularity - md1.Popularity;
    if (graph[url2].length == graph[url1].length)
      return url1.localeCompare(url2); // they are equal in the order; sort lexicographically for determinism
    return graph[url2].length - graph[url1].length;
  });
  let nodeToRow = {};
  withMetadata.forEach(([url, _], i) => {
    nodeToRow[url] = i;
  });
  return nodeToRow;
}

function getNodeCol(depth, i) {
  if (depth == 0)
    return COL_WIDTH - 1;
  return (depth * COL_WIDTH) + (i % COL_WIDTH);
}

// Hack for allowing us to test some of these functions without setting up a
// proper webpack build pipeline.
if (typeof exports == 'undefined') {
  window.addEventListener('load', (_) => {
    renderGraph();
  });
} else {
  exports.displayUrl = displayUrl;
  exports.getNodeRows = getNodeRows;
}
