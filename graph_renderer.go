// GraphParameters is a class
// GraphParameters has a public Rect, ScreenBounds
// It has a public vector2, the yvalues minmax
// It has a public list of list of float, called yvalueslist
// and it has a public list of float, called xvalues
// it has an ienumerator of colors called Cols
//
// it has a single method, GetYValuesMinMax that scans the YvaluesList,
// and returns a min and max over the whole list of lists
//
// GraphRenderer is a class, which is a MonoBehavor
//
// GraphRenderer has a static instance of GraphRenderer
// RenderGraph takes a GraphParameters and enqueues it in a member, graphs
// graphs is a queue of graphParameters
// graph_mat is a member which is a material - a shader material?
//
// Awake is a method, which stores "this" into Instance
//
// OnPostRender is a method, which empties the queue of graphs, and calls PostRenderGraph on each
//
// PostRenderGraph takes a GraphParameters and does a lot of GL drawing

package main
