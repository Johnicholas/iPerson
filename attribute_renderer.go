// AttributeRenderer graphs all attributes over time, on screen
//
// AttributeRenderer has a float, sample interval, in [0, 2]
//
// has a history size, an int, initially 50
//
// has a screen width is a float, initially 1, range [0, 1]
//
// has a screen height is a float, initially 1, range [0, 1]
//
// has a screen y prop is a float, initially 1, range [0, 1]
//
// has a bars_originprop, which is a vector2
// has a bars_size, which is a vector2
//
// has a slider, a float, with range [0, 1]
//
// has attribute_histories, a map from person_attributes to queue of floats
//
// has action_score_histories, a map from game object to queue of floats
//
// has actions_history, a queue of string
//
// has current history size, an int
//
// Source is either ATTRIBUTES or SCORES
// has a current source, which is initially ATTRIBUTES
// has pending_action, which is a string
//
// has a reference to the DecisionFlex
//
// has a boolean, is_paused
//
// implements a method Awake() void
// Awake sets IsPaused to isPaused?
// it gets the components of type PersonAttribute, and stores them in a local, attributes
// it creates a new empty map from person attribute to queue of float
// it iterates over all the attributes, and initializes the map with that key
// to point to an empty queue
//
// implements a method Start() void
// it finds a DecisionFlex object, and stores it in a member
// if it found it it also subscribes to OnNewAction
// it gets the components of type Action, and stores them in a local, allActions
// for each action in allActions, it initializes the action score histories dictionary
// to point to a new queue of float
//
// Then it calls StartCoroutine(Sampler())
//
// IsPaused is a public bool?
//
// Sampler returns an IEnumerator
// it runs three PerformAction()/RecordActionScores()/RecordAttributes() triples to "warm up"
// forever, it invokes RecordActionScores()/RecordAttributes()
// and then "yields" WaitForSeconds(sampleInterval)
//
// RecordAttributes is a method that
// for each attribute history in the member attributeHistories,
// we get the current value of the key (which is the attribute)
// and do a PushWithRestrictedSize
// finally, we PushWithRestrictedSize the pendingAction into actionsHistory,
// and null the pendingAction
//
// RecordActionScores is a method that
// pulls AllLastSelections from the decisionflex object
// sets isActionTurn according to whether pendingAction isNullOrEmpty
// for each actionScore,
//   gets the actionObject of that actionScore
//   gets the recentActionScores queue from the actionScoreHistories
//   assigns scoreValue to be the actionScore's score, or zero if it is actionTurn
//   and pushWithRestrictedSize the score value into the queue
//
// PushWithRestrictedSize enqueues an item, and if the
// size of the queue is larger than a limit, also dequeues an item
//
// OnGUI is a method that:
// 1. delegates to RenderAttributeBars
// 2. constructs a Rect from padding, screen width and ehight, and width and height props
// 3. delegates to RenderControlAroundGraph passing that rect
// 4. delegates to RenderHistories
// (specifically either attributeHistories or actionScoreHistories, depending on currentSource)
// 5. delegates to RenderActionsOnGraph
//
// RenderAttributeBars is a method that:
// 1. builds a new BarRect
// 2. for each attribute in attributeHistories.keys
//       calls GUI and GUILayout to build a bar chart based on attr.value
//
// RenderControlAroundGraph is a method that
// is a sequence of calls to GUILayout
//
// RenderHistories is a method that takes a histories, graphParameters, and toString
// for each value in histories.values, copies the value into a yValuesList
// delegates to GraphRenderer.RenderGraph
// for each history in histories
// uses LastNotableValueIn and CalculateScreenCoord to
// Label the history with a GUI.Label
//
// LastNotableValueIn scans the list of values for the last positive one,
// or returns the last one if they are all nonpositive

// RenderActionsOnGraph is a method that:
// for each action in actionsHistory
// calls GUI.Label to label the action
// using CalculateScreenCoord to decide where
//
// OnAction gets the winningAction, converts it to a string and stores it in pendingAction
//
// CalculateScreenCoord returns an appropriate place to put something
//
// ColorGenerator yields a sequence of random colors
package main
