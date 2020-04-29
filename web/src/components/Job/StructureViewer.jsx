import LiteMol from "litemol";
import "litemol/dist/css/LiteMol-plugin.css";
import React from "react";
import "../../styles/components/structure-viewer.scss";
import { ResultsContext } from "./ResultsContext";

const Transformer = LiteMol.Bootstrap.Entity.Transformer;
const Transform = LiteMol.Bootstrap.Tree.Transform;

export default class StructureViewer extends React.Component {
  constructor(props) {
    super(props);
    this.state = { plugin: {} };
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

  load() {
    this.clear();

    let id = this.context.results.PDB.ID;
    let action = Transform.build()
      .add(this.state.plugin.context.tree.root, Transformer.Data.Download, {
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
        { style: surfaceStyle },
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
        { style: surfaceStyle },
        { isHidden: true, ref: "surface" }
      );

    this.state.plugin.applyTransform(action).then(() => {
      this.applyTheme("polymer-visual", this.createTheme());
      this.applyTheme("surface", this.createTheme(0.2, 0x0d6273));
      this.applyTheme("surface-het", this.createTheme(0.4, 0x00fffb));
    });

    this.state.plugin.command(LiteMol.Bootstrap.Command.Layout.SetState, {
      collapsedControlsLayout:
        LiteMol.Bootstrap.Components.CollapsedControlsLayout.Landscape,
      hideControls: true,
    });
  }

  applyTheme(ref, theme) {
    var plugin = this.state.plugin;

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
    var plugin = this.state.plugin;
    let model = plugin.context.select("model")[0];

    const fallbackColor = LiteMol.Visualization.Color.fromHex(uniform);
    const selectionColor = LiteMol.Visualization.Color.fromHex(0xf15a29);
    const highlightColor = LiteMol.Visualization.Color.fromHex(0xff752b);
    const mutedColor = LiteMol.Visualization.Color.fromHex(0x163d40);

    let colors = new Map();

    let SIFTSUnp = this.context.results.PDB.SIFTS.UniProt;
    Object.keys(SIFTSUnp)
      .filter((k) => {
        return this.context.results.UniProt.ID != k;
      })
      .forEach((id) => {
        SIFTSUnp[id].mappings.forEach((chain) => {
          colors.set(chain.chain_id, mutedColor);
        });
      });

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

  highlightHet(id) {
    this.clearHighlight();
    var plugin = this.state.plugin;
    var model = plugin.context.select("model")[0];
    let query = LiteMol.Core.Structure.Query.residuesByName(id);
    LiteMol.Bootstrap.Command.Molecule.Highlight.dispatch(plugin.context, {
      model,
      query,
      isOn: true,
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

  selectFocus(chain, start, end) {
    this.focus(chain, start, end);
    this.highlightResidues([
      { Chain: chain, Position: start, PositionEnd: end },
    ]);
    if (start - end == 0) {
      this.select(chain, start, end);
    }
  }

  render() {
    return <div id="litemol" />;
  }
}
StructureViewer.contextType = ResultsContext;
