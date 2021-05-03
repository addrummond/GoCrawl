import { expect } from 'chai'
import { displayUrl, getNodeRows } from './render.js' 

describe('displayUrl', () => {
  it('yields empty string if stripping empty string from empty string', () => {
    expect(displayUrl("", "")).to.equal("");
  });
  it('yields the original string with preceding / if the prefix is empty', () => {
    expect(displayUrl("", "foo")).to.equal("/foo");
    expect(displayUrl("", "foo/bar")).to.equal("/foo/bar");
  });
  it('yields the original URL if it is identical to the prefix', () => {
    expect(displayUrl("http://google.com", "http://google.com")).to.equal("http://google.com");
    expect(displayUrl("http://google.com/", "http://google.com/")).to.equal("http://google.com/");
  });
  it('strips the prefix from a URL', () => {
    expect(displayUrl("http://google.com", "http://google.com/dir/page")).to.equal("/dir/page");
    expect(displayUrl("http://google.com/", "http://google.com/dir/page")).to.equal("/dir/page");
  });
});

describe('getNodeRows', () => {
  it('orders nodes by (number of outgoing edges DESC, popularity DESC)', () => {
    // The outgoing edges will never be followed, just counted, so we don't need
    // to construct a full graph.
    let graph = {};
    let metadata = {};
    graph.A = [null, null, null];
    graph.B = [null];
    graph.C = [null, null];
    graph.D = [null];
    graph.E = [null, null];
    graph.F = [null];

    metadata.A = {Popularity: 4};
    metadata.B = {Popularity: 3};
    metadata.C = {Popularity: 3};
    metadata.D = {Popularity: 9};
    metadata.E = {Popularity: 9};
    metadata.F = {Popularity: 0};

    const expectedOrder = {A: 0, E: 1, C: 2, D: 3, B: 4, F: 5};
    const rows = getNodeRows(graph, metadata);
    expect(rows).to.deep.equal(expectedOrder);
  });
});
