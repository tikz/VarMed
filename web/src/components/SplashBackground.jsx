import React from "react";
import LiteMol from "litemol";
import "litemol/dist/css/LiteMol-plugin.css";
import "../styles/components/structure-viewer.scss";

const Transformer = LiteMol.Bootstrap.Entity.Transformer;
const Transform = LiteMol.Bootstrap.Tree.Transform;

function mousemoveListener(e) {
  if (e.isTrusted) {
    e.stopPropagation();
  }
}

function simulate(element, eventName) {
  var options = extend(defaultOptions, arguments[2] || {});
  var oEvent,
    eventType = null;

  for (var name in eventMatchers) {
    if (eventMatchers[name].test(eventName)) {
      eventType = name;
      break;
    }
  }

  if (!eventType)
    throw new SyntaxError(
      "Only HTMLEvents and MouseEvents interfaces are supported"
    );

  if (document.createEvent) {
    oEvent = document.createEvent(eventType);
    if (eventType == "HTMLEvents") {
      oEvent.initEvent(eventName, options.bubbles, options.cancelable);
    } else {
      oEvent.initMouseEvent(
        eventName,
        options.bubbles,
        options.cancelable,
        document.defaultView,
        options.button,
        options.pointerX,
        options.pointerY,
        options.pointerX,
        options.pointerY,
        options.ctrlKey,
        options.altKey,
        options.shiftKey,
        options.metaKey,
        options.button,
        element
      );
    }
    element.dispatchEvent(oEvent);
  } else {
    options.clientX = options.pointerX;
    options.clientY = options.pointerY;
    var evt = document.createEventObject();
    oEvent = extend(evt, options);
    element.fireEvent("on" + eventName, oEvent);
  }
  return element;
}

function extend(destination, source) {
  for (var property in source) destination[property] = source[property];
  return destination;
}

var eventMatchers = {
  HTMLEvents: /^(?:load|unload|abort|error|select|change|submit|reset|focus|blur|resize|scroll)$/,
  MouseEvents: /^(?:click|dblclick|mouse(?:down|up|over|move|out))$/,
};

var defaultOptions = {
  pointerX: 0,
  pointerY: 0,
  button: 0,
  ctrlKey: false,
  altKey: false,
  shiftKey: false,
  metaKey: false,
  bubbles: true,
  cancelable: true,
};

export default class SplashBackground extends React.Component {
  constructor(props) {
    super(props);
    this.litemolRef = React.createRef();
    this.state = { plugin: {}, collapsed: false };
  }

  componentDidMount() {
    this.plugin = LiteMol.Plugin.create({
      target: "#splash",
      viewportBackground: "#1c1e20",
      layoutState: {
        hideControls: true,
        isExpanded: false,
      },
    });

    document.addEventListener("mousemove", mousemoveListener, {
      capture: true,
    });

    this.load();
    setInterval(function () {
      simulate(document.querySelector("#splash canvas"), "mousedown", {
        pointerX: 100,
        pointerY: 100,
      });
      simulate(document.querySelector("#splash canvas"), "mousemove", {
        pointerX: 99,
        pointerY: 99,
      });
      simulate(document.querySelector("#splash canvas"), "mouseup", {
        pointerX: 100,
        pointerY: 100,
      });
    }, 50);

    const that = this;
    let i = 1;

    setInterval(function () {
      let h = setInterval(function () {
        that.highlight(i, 5);
        i = i < 1456 ? i + 1 : 1;
        if (i < 1456) {
          i++;
        } else {
          i = 1;
          clearInterval(h);
        }
      }, 10);
    }, 60000);
  }

  load() {
    this.clear();

    let id = "5MZO";
    let action = Transform.build()
      .add(this.plugin.context.tree.root, Transformer.Data.Download, {
        url: API_URL + `/api/structure/cif/${id}`,
        type: "String",
        id,
      })
      .then(Transformer.Data.ParseCif, { id }, { isBinding: true })
      .then(
        Transformer.Molecule.CreateFromMmCif,
        { blockIndex: 0 },
        { isBinding: true }
      )
      .then(
        Transformer.Molecule.CreateModel,
        { modelIndex: 0 },
        { isBinding: false, ref: "model" }
      );

    let sel = action;
    sel.then(Transformer.Molecule.CreateMacromoleculeVisual, {
      polymer: true,
      polymerRef: "polymer-visual",
      het: true,
      water: true,
    });

    let surfaceStyle = {
      type: "Surface",
      params: {
        probeRadius: 0.4,
        density: 2,
        smoothing: 3,
        isWireframe: false,
      },
      theme: {
        template:
          LiteMol.Bootstrap.Visualization.Molecule.Default.UniformThemeTemplate,
        transparency: { alpha: 0.2 },
      },
    };

    sel
      .then(
        Transformer.Molecule.CreateSelectionFromQuery,
        {
          query: LiteMol.Core.Structure.Query.hetGroups(),
          name: "Het",
          silent: true,
        },
        {}
      )
      .then(
        Transformer.Molecule.CreateVisual,
        { style: surfaceStyle },
        { isHidden: true, ref: "surface-het" }
      );

    sel
      .then(
        Transformer.Molecule.CreateSelectionFromQuery,
        {
          query: LiteMol.Core.Structure.Query.nonHetPolymer(),
          name: "Surface",
          silent: true,
        },
        {}
      )
      .then(
        Transformer.Molecule.CreateVisual,
        { style: surfaceStyle },
        { isHidden: true, ref: "surface" }
      );

    this.plugin.applyTransform(action).then(() => {
      this.applyTheme("polymer-visual", this.createTheme());
      this.applyTheme("surface", this.createTheme(0.2, 0x0d6273));
      this.applyTheme("surface-het", this.createTheme(0.4, 0x00fffb));

      document.getElementById("splash").classList.add("started");
      this.focus("A", 350, 400);
    });

    this.plugin.command(LiteMol.Bootstrap.Command.Layout.SetState, {
      collapsedControlsLayout:
        LiteMol.Bootstrap.Components.CollapsedControlsLayout.Portrait,
      hideControls: true,
    });

    this.plugin.context.scene.scene.updateOptions({
      fogFactor: 0.01,
      cameraSpeed: 0.5,
    });
  }

  applyTheme(ref, theme) {
    var plugin = this.plugin;

    const visuals = plugin.selectEntities(
      LiteMol.Bootstrap.Tree.Selection.byRef(ref)
        .subtree()
        .ofType(LiteMol.Bootstrap.Entity.Molecule.Visual)
    );

    for (const v of visuals) {
      plugin.command(LiteMol.Bootstrap.Command.Visual.UpdateBasicTheme, {
        visual: v,
        theme,
      });
    }
  }

  createTheme(alpha = 1, uniform = 0x006870) {
    var plugin = this.plugin;
    let model = plugin.context.select("model")[0];

    const fallbackColor = LiteMol.Visualization.Color.fromHex(uniform);
    const selectionColor = LiteMol.Visualization.Color.fromHex(0xf15a29);
    const highlightColor = LiteMol.Visualization.Color.fromHex(0xff752b);

    let colors = new Map();

    colors.set("Uniform", fallbackColor);
    colors.set("Selection", selectionColor);
    colors.set("Highlight", highlightColor);

    let theme = LiteMol.Bootstrap.Visualization.Molecule.createColorMapThemeProvider(
      (m) => ({
        index: m.data.atoms.chainIndex,
        property: m.data.chains.asymId,
      }),
      colors,
      fallbackColor
    )(model, {
      colors: colors,
      transparency: { alpha: alpha },
      isSticky: true,
    });
    return theme;
  }

  showSurface(visible) {
    this.setVisibility("surface", visible);
    this.setVisibility("surface-het", visible);
  }

  setVisibility(ref, visible) {
    let entity = this.plugin.context.select(ref)[0];

    LiteMol.Bootstrap.Command.Entity.SetVisibility.dispatch(
      this.plugin.context,
      { entity, visible: visible }
    );
  }

  clear() {
    LiteMol.Bootstrap.Command.Tree.RemoveNode.dispatch(
      this.plugin.context,
      this.plugin.context.tree.root
    );
  }

  highlight(pos, q) {
    this.clearHighlight();
    var plugin = this.plugin;
    var model = plugin.context.select("model")[0];

    for (let i = pos; i < pos + q; i++) {
      let query = LiteMol.Core.Structure.Query.sequence(
        null,
        "A",
        { seqNumber: i },
        {
          seqNumber: i,
        }
      );
      LiteMol.Bootstrap.Command.Molecule.Highlight.dispatch(plugin.context, {
        model,
        query,
        isOn: true,
      });
    }
  }

  clearHighlight() {
    var plugin = this.plugin;
    var model = plugin.context.select("model")[0];
    let query = LiteMol.Core.Structure.Query.everything();
    LiteMol.Bootstrap.Command.Molecule.Highlight.dispatch(plugin.context, {
      model,
      query,
      isOn: false,
    });
  }

  focus(chain, start, end) {
    var plugin = this.plugin;
    var model = plugin.context.select("model")[0];
    let query = LiteMol.Core.Structure.Query.sequence(
      null,
      chain,
      { seqNumber: start },
      { seqNumber: end }
    );
    LiteMol.Bootstrap.Command.Molecule.FocusQuery.dispatch(plugin.context, {
      model,
      query,
      isOn: true,
    });
  }

  select(chain, start, end) {
    var plugin = this.plugin;
    var model = plugin.context.select("model")[0];
    let query = LiteMol.Core.Structure.Query.sequence(
      null,
      chain,
      { seqNumber: start },
      { seqNumber: end }
    );
    LiteMol.Bootstrap.Command.Molecule.CreateSelectInteraction.dispatch(
      plugin.context,
      { entity: model, query: query }
    );
  }

  selectFocus(chain, start, end) {
    this.select(chain, start, end);
    this.focus(chain, start, end);
    console.log("blah");
  }

  render() {
    return <div id="splash" />;
  }
}
