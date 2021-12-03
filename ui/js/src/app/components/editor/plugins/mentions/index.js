'use strict';

Object.defineProperty(exports, '__esModule', { value: true });

var prosemirrorState = require('prosemirror-state');
var prosemirrorView = require('prosemirror-view');

/**
 *
 * @param {String} mentionTrigger
 * @param {String} hashtagTrigger
 * @param {bool} allowSpace
 * @returns {Object}
 */


function getRegexp(mentionTrigger, hashtagTrigger, emojiTrigger, allowSpace) {
let spaceReg = new RegExp("(^|\\s)" + mentionTrigger + "([\\w-\\+]+\\s?[\\w-\\+]*)$")
let nonSpaceReg = new RegExp("(^|\\s)" + mentionTrigger + "([\\w-/\\+]+)$")
  var mention = allowSpace ? spaceReg : nonSpaceReg ;

  // hashtags should never allow spaces. I mean, what's the point of allowing spaces in hashtags?
  var tag = new RegExp("(^|\\s)" + hashtagTrigger + "([\\w-]+)$");

  var emoji = new RegExp("(^|\\s)" + emojiTrigger + "([\\w-]+)$");

  return {
    mention: mention,
    tag: tag,
    emoji: emoji
  };
}

/**
 *
 * @param {ResolvedPosition} $position https://prosemirror.net/docs/ref/#model.Resolved_Positions
 * @param {JSONObject} opts
 * @returns {JSONObject}
 */
function getMatch($position, opts) {
  // take current para text content upto cursor start.
  // this makes the regex simpler and parsing the matches easier.
  var parastart = $position.before();
  const text = $position.doc.textBetween(parastart, $position.pos, "\n", "\0");

  var regex = getRegexp(opts.mentionTrigger, opts.hashtagTrigger, opts.emojiTrigger, opts.allowSpace);

  // only one of the below matches will be true.
  var mentionMatch = text.match(regex.mention);
  var tagMatch = text.match(regex.tag);
  var emojiMatch = text.match(regex.emoji);

  var match = mentionMatch || tagMatch || emojiMatch;

  // set type of match
  var type;
  if (mentionMatch) {
    type = "mention";
  } else if (tagMatch) {
    type = "tag";
  } else if (emojiMatch) {
    type = "emoji";
  }

  // if match found, return match with useful information.
  if (match) {
    // adjust match.index to remove the matched extra space
    match.index = match[0].startsWith(" ") ? match.index + 1 : match.index;
    match[0] = match[0].startsWith(" ") ? match[0].substring(1, match[0].length) : match[0];

    // The absolute position of the match in the document
    var from = $position.start() + match.index;
    var to = from + match[0].length;

    var queryText = match[2];

    return {
      range: { from: from, to: to },
      queryText: queryText,
      type: type
    };
  }
  // else if no match don't return anything.
}

/**
 * Util to debounce call to a function.
 * >>> debounce(function(){}, 1000, this)
 */
const debounce = function () {
  var timeoutId = null;
  return function (func, timeout, context) {
    context = context || this;
    clearTimeout(timeoutId);
    timeoutId = setTimeout(function () {
      func.apply(context, arguments);
    }, timeout);

    return timeoutId;
  };
}();

var getNewState = function () {
  return {
    active: false,
    range: {
      from: 0,
      to: 0
    },
    type: "", //mention or tag
    text: "",
    suggestions: [],
    index: 0 // current active suggestion index
  };
};

/**
 * @param {JSONObject} opts
 * @returns {Plugin}
 */
function getMentionsPlugin(opts) {
  // default options
  var defaultOpts = {
    mentionTrigger: "@",
    hashtagTrigger: "#",
    emojiTrigger: ":",
    allowSpace: false,
    getSuggestions: (type, text, cb) => {
      cb([]);
    },
    activeClass: "suggestion-item-active",
    suggestionTextClass: "suggestion",
    maxNoOfSuggestions: 13,
    delay: 200
  };

  var opts = Object.assign({}, defaultOpts, opts);

  // timeoutId for clearing debounced calls
  var showListTimeoutId = null;

  // dropdown element
  var el = document.createElement("div");

  // ----- methods operating on above properties -----

  var showList = function (view, state, suggestions, opts) {
    el.innerHTML = opts.getSuggestionsHTML(suggestions, state.type, state.text);

    // attach new item event handlers
    el.querySelectorAll(".suggestion-item").forEach(function (itemNode, index) {
      itemNode.addEventListener("click", function () {
        select(view, state, opts);
        view.focus();
      });
      // TODO: setIndex() needlessly queries.
      // We already have the itemNode. SHOULD OPTIMIZE.
      itemNode.addEventListener("mouseover", function () {
        setIndex(index, state, opts);
      });
      itemNode.addEventListener("mouseout", function () {
        setIndex(index, state, opts);
      });
    });

    // highlight first element by default - like Facebook.
    addClassAtIndex(state.index, opts.activeClass);

    // get current @mention span left and top.
    // TODO: knock off domAtPos usage. It's not documented and is not officially a public API.
    // It's used currently, only to optimize the the query for textDOM

    // TODO: should add null check case for textDOM

    // TODO: think about outsourcing this positioning logic as options
    let ch = document.querySelector(`#${opts.id}`)
    ch.appendChild(el);
    el.className = "suggestion-container"
    el.style.zIndex = 50000;
    el.style.position = "absolute";
    el.style.left = `1rem`
    el.style.paddingRight = `2rem`
    el.style.width = `100%`

    el.style.bottom = `76px`
    el.style.display = "block";
    el.setAttribute("data-active", true)
  };

  var hideList = function () {
    el.style.display = "none";
    el.setAttribute("data-active", false)
  };

  var removeClassAtIndex = function (index, className) {
    var itemList = el.querySelector(".suggestion-item-list").childNodes;
    var prevItem = itemList[index];
    prevItem.classList.remove(className);
  };

  var addClassAtIndex = function (index, className) {
    var itemList = el.querySelector(".suggestion-item-list").childNodes;
    var prevItem = itemList[index];
    if(prevItem) {
      prevItem.classList.add(className);
    }
  };

  var setIndex = function (index, state, opts) {
    removeClassAtIndex(state.index, opts.activeClass);
    state.index = index;
    addClassAtIndex(state.index, opts.activeClass);
  };

  var goNext = function (view, state, opts) {
    removeClassAtIndex(state.index, opts.activeClass);
    state.index++;
    state.index = state.index === state.suggestions.length ? 0 : state.index;
    addClassAtIndex(state.index, opts.activeClass);
  };

  var goPrev = function (view, state, opts) {
    removeClassAtIndex(state.index, opts.activeClass);
    state.index--;
    state.index = state.index === -1 ? state.suggestions.length - 1 : state.index;
    addClassAtIndex(state.index, opts.activeClass);
  };

  var select = function (view, state, opts) {
    var item = state.suggestions[state.index];
    var attrs;
    if (state.type === "mention") {
      /*
      if(item == null) {
        return
      }
      */
      attrs = {
        username: item?.username,
        user_id: item?.user_id,
      };
    } else if (state.type == "tag"){
      attrs = {
        name: item?.name,
        path: item?.path,
      };
    }

    let from = state.range.from
    let to = state.range.to
    var tr;
    if (state.type === "mention" && (!opts.key || opts.key != 'space')) {
      var node = view.state.schema.nodes[state.type].create(attrs);
      tr = view.state.tr.replaceWith(from, to, node);
      tr = tr.insertText(" ")
      window.insertedMention = true
      setTimeout(() => {
        window.insertedMention = false
      }, 100)
    }
    if (state.type === "tag") {
      var node = view.state.schema.nodes[state.type].create(attrs);
      tr = view.state.tr.replaceWith(from, to, node);
      tr = tr.insertText(" ")
      window.insertedTag = true
      setTimeout(() => {
        window.insertedTag = false
      }, 100)
    }
    if(state.type =="emoji") {
      var node = view.state.schema.text(`${item.unicode} `)
      tr = view.state.tr
      tr = view.state.tr.replaceWith(from, to, node);
      window.insertedEmoji = true
      setTimeout(() => {
        window.insertedEmoji = false
      }, 100)
    }

    //var newState = view.state.apply(tr);
    //view.updateState(newState);
    view.dispatch(tr)
  };

  /**
   * See https://prosemirror.net/docs/ref/#state.Plugin_System
   * for the plugin properties spec.
   */
  return new prosemirrorState.Plugin({
    key: new prosemirrorState.PluginKey("autosuggestions"),

    // we will need state to track if suggestion dropdown is currently active or not
    state: {
      init() {
        return getNewState();
      },

      apply(tr, state) {
        // compute state.active for current transaction and return
        var newState = getNewState();
        var selection = tr.selection;
        if (selection.from !== selection.to) {
          return newState;
        }

        const $position = selection.$from;
        const match = getMatch($position, opts);

        // if match found update state
        if (match) {
          newState.active = true;
          newState.range = match.range;
          newState.type = match.type;
          newState.text = match.queryText;
        }

        return newState;
      }
    },

    // We'll need props to hi-jack keydown/keyup & enter events when suggestion dropdown
    // is active.
    props: {
      handleKeyDown(view, e) {
        var state = this.getState(view.state);

        // don't handle if no suggestions or not in active mode
        if (!state.active && !state.suggestions.length) {
          return false;
        }

        // if any of the below keys, override with custom handlers.
        var down, up, enter, esc, space, tab;
        enter = e.keyCode === 13;
        down = e.keyCode === 40;
        up = e.keyCode === 38;
        esc = e.keyCode === 27;
        space = e.keyCode === 32;
        tab = e.keyCode === 9;

        if(space) {
          /*
          opts.key = 'space'
          clearTimeout(showListTimeoutId);
          hideList();
          this.state = getNewState();
          select(view, state, opts);
          return false
          */
            return false
        }

        if (down || tab) {
          goNext(view, state, opts);
          return true;
        } else if (up) {
          goPrev(view, state, opts);
          return true;
        } else if (enter) {
          select(view, state, opts);
          hideList();
          this.state = getNewState();
          return true;
        } else if (esc) {
          clearTimeout(showListTimeoutId);
          hideList();
          this.state = getNewState();
          return true;
        } else {
          // didn't handle. handover to prosemirror for handling.
          return false;
        }
      },

      // to decorate the currently active @mention text in ui
      decorations(editorState) {
        const { active, range } = this.getState(editorState);
        var state = this.getState(editorState);

        if (!active) return null;

        let options = {
          nodeName: "span",
          class: opts.suggestionTextClass
        }

        if(state.type == 'tag') {
          options.class += " hashtag"
          options["data-hashtag"] = state.text
          options["data-focusable"] = true
        }

        return prosemirrorView.DecorationSet.create(editorState.doc, [prosemirrorView.Decoration.inline(range.from, range.to, options)]);
      }
    },

    // To track down state mutations and add dropdown reactions
    view() {
      return {
        update: view => {
          var state = this.key.getState(view.state);
          if (!state.text) {
            hideList();
            clearTimeout(showListTimeoutId);
            return;
          }
          // debounce the call to avoid multiple requests
          showListTimeoutId = debounce(function () {
            // get suggestions and set new state
            opts.getSuggestions(state.type, state.text, function (suggestions) {
              // update `state` argument with suggestions
              if(state.type === 'tag') {
                //suggestions.unshift({tag: state.text})
              }
              state.suggestions = suggestions;
              showList(view, state, suggestions, opts);
            });
          }, opts.delay, this);
        }
      };
    }
  });
}

/**
 * See https://prosemirror.net/docs/ref/#model.NodeSpec
 */
const mentionNode = {
  group: "inline",
  inline: true,
  atom: true,

  attrs: {
    username: "",
    user_id: "",
  },

  selectable: false,
  draggable: false,

  toDOM: node => {
    return ["span", {
      class: "mention",
      "data-focusable": true,
      "data-username": node.attrs.username,
      "data-userid": node.attrs.user_id,
    },
      "@" + node.attrs.username
    ];
  },

  parseDOM: [{
    // match tag with following CSS Selector
    tag: "span[data-username]",

    getAttrs: dom => {
      var username = dom.getAttribute("data-username");
      var user_id = dom.getAttribute("data-userid");
      return {
        username: username,
        user_id: user_id,
      };
    }
  }]
};

/**
 * See https://prosemirror.net/docs/ref/#model.NodeSpec
 */
const tagNode = {
  group: "inline",
  inline: true,
  atom: true,

  attrs: {
    name: "",
    path: "",
  },

  selectable: false,
  draggable: false,

  toDOM: node => {
    return ["span", {
      "data-path": node.attrs.path,
      "data-name": node.attrs.name,
      class: "tag"
    }, 
      "#" + node.attrs.name 
    ];
  },

  parseDOM: [{
    // match tag with following CSS Selector
    tag: "span[data-name]",

    getAttrs: dom => {
      var name = dom.getAttribute("data-name");
      var path = dom.getAttribute("data-path");
      return {
        name: name,
        path: path
      };
    }
  }]
};

/**
 *
 * @param {OrderedMap} nodes
 * @returns {OrderedMap}
 */
function addMentionNodes(nodes) {
  return nodes.append({
    mention: mentionNode
  });
}

/**
 *
 * @param {OrderedMap} nodes
 * @returns {OrderedMap}
 */
function addTagNodes(nodes) {
  return nodes.append({
    tag: tagNode
  });
}

exports.getMentionsPlugin = getMentionsPlugin;
exports.addMentionNodes = addMentionNodes;
exports.addTagNodes = addTagNodes;
exports.tagNode = tagNode;
exports.mentionNode = mentionNode;
