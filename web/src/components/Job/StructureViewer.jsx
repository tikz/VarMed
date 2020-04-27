import React from "react";
import LiteMol from "litemol";
import "litemol/dist/css/LiteMol-plugin.css";
import "../../styles/components/structure-viewer.scss";

let Transformer = LiteMol.Bootstrap.Entity.Transformer;
let Transform = LiteMol.Bootstrap.Tree.Transform;

export default class StructureViewer extends React.Component {
  constructor(props) {
    super(props);
    this.state = { plugin: {}, res: {} };
  }

  componentDidMount() {
    this.setState({
      plugin: LiteMol.Plugin.create({
        target: "#litemol",
        viewportBackground: "#1c1e20",
        layoutState: {
          hideControls: true,
          isExpanded: false,
        },
      }),
    });
  }

  load(res) {
    this.clear();
    this.setState({ res: res });

    let surfaceColors = LiteMol.Bootstrap.Immutable.Map()
      .set("Uniform", LiteMol.Visualization.Color.fromHex(0x0d6273))
      .set("Selection", LiteMol.Visualization.Color.fromHex(0xf15a29))
      .set("Highlight", LiteMol.Visualization.Color.fromHex(0xff752b));

    let hetColors = surfaceColors.set(
      "Uniform",
      LiteMol.Visualization.Color.fromHex(0x00fffb)
    );
    let polymerSurfaceStyle = {
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
        colors: surfaceColors,
      },
    };

    let hetSurfaceStyle = {
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
        transparency: { alpha: 0.4 },
        colors: hetColors,
      },
    };

    let id = res.PDB.ID;
    let action = Transform.build()
      .add(this.state.plugin.context.tree.root, Transformer.Data.Download, {
        url: `http://localhost:8888/api/structure/cif/${id}`,
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

    sel = action;
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
        { style: hetSurfaceStyle },
        { isHidden: true, ref: "surface-het" }
      );

    sel = action;
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
        { style: polymerSurfaceStyle },
        { isHidden: true, ref: "surface" }
      );

    this.state.plugin.applyTransform(action).then(() => {
      this.applyTheme();
    });
  }

  showSurface(visible) {
    this.setVisibility("surface", visible);
    this.setVisibility("surface-het", visible);
  }

  setVisibility(ref, visible) {
    let entity = this.state.plugin.context.select(ref)[0];

    LiteMol.Bootstrap.Command.Entity.SetVisibility.dispatch(
      this.state.plugin.context,
      { entity, visible: visible }
    );
  }

  clear() {
    LiteMol.Bootstrap.Command.Tree.RemoveNode.dispatch(
      this.state.plugin.context,
      this.state.plugin.context.tree.root
    );
  }

  highlight(chain, start, end) {
    this.clearHighlight();
    var plugin = this.state.plugin;
    var model = plugin.context.select("model")[0];
    let query = LiteMol.Core.Structure.Query.sequence(
      null,
      chain,
      { seqNumber: start },
      { seqNumber: end }
    );
    LiteMol.Bootstrap.Command.Molecule.Highlight.dispatch(plugin.context, {
      model,
      query,
      isOn: true,
    });
  }

  highlightResidues(residues) {
    this.clearHighlight();
    var plugin = this.state.plugin;
    var model = plugin.context.select("model")[0];
    residues.forEach((r) => {
      let query = LiteMol.Core.Structure.Query.sequence(
        null,
        r.Chain,
        { seqNumber: r.Position },
        {
          seqNumber: r.PositionEnd !== undefined ? r.PositionEnd : r.Position,
        }
      );
      LiteMol.Bootstrap.Command.Molecule.Highlight.dispatch(plugin.context, {
        model,
        query,
        isOn: true,
      });
    });
  }

  clearHighlight() {
    var plugin = this.state.plugin;
    var model = plugin.context.select("model")[0];
    let query = LiteMol.Core.Structure.Query.everything();
    LiteMol.Bootstrap.Command.Molecule.Highlight.dispatch(plugin.context, {
      model,
      query,
      isOn: false,
    });
  }

  focus(chain, start, end) {
    var plugin = this.state.plugin;
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
    var plugin = this.state.plugin;
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

  applyTheme() {
    var plugin = this.state.plugin;
    let colors = LiteMol.Core.Utils.FastMap.create();
    colors.set("Uniform", LiteMol.Visualization.Color.fromHex(0x006870));
    colors.set("Selection", LiteMol.Visualization.Color.fromHex(0xf15a29));
    colors.set("Highlight", LiteMol.Visualization.Color.fromHex(0xff752b));
    let theme = LiteMol.Bootstrap.Visualization.Molecule.uniformThemeProvider(
      void 0,
      { colors }
    );

    const visuals = plugin.selectEntities(
      LiteMol.Bootstrap.Tree.Selection.byRef("polymer-visual")
        .subtree()
        .ofType(LiteMol.Bootstrap.Entity.Molecule.Visual)
    );
    for (const v of visuals) {
      plugin.command(LiteMol.Bootstrap.Command.Visual.UpdateBasicTheme, {
        visual: v,
        theme,
      });
    }
    plugin.command(LiteMol.Bootstrap.Command.Layout.SetState, {
      collapsedControlsLayout:
        LiteMol.Bootstrap.Components.CollapsedControlsLayout.Landscape,
      hideControls: true,
    });
  }

  render() {
    return <div id="litemol" />;
  }
}
