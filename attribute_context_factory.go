package main

// pushes the iPerson attributes into a single IContext
//
// an AttributeContextFactory is a SingleContextFactory
//
// AttributeContextFactory has a list of all attributes
//
// AttributeContextFactory implements SingleContext by
// creating a new ContextDictionary, and loops over all attributes,
// setting name, value pairs into the context
// and returning it
//
// AttributeContextFactory implements Awake by searching for all
// components of type PersonAttribute in the children of the parent
