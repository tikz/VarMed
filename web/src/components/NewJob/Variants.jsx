import {
  Box,
  Checkbox,
  FormControlLabel,
  Typography,
  Link,
} from "@material-ui/core";
import React from "react";
import VariantInput from "./VariantInput";
import ChipArray from "./ChipArray";

export class Variants extends React.Component {
  constructor(props) {
    super(props);
    this.handleChange = this.handleChange.bind(this);
    this.handleDelete = this.handleDelete.bind(this);
  }

  handleChange(e) {
    this.props.setAnnotated(e.target.checked);
  }

  handleDelete(chip) {
    this.props.setVariants(
      this.props.variants.filter((c) => c.key !== chip.key)
    );
  }

  render() {
    let unpSeqURL =
      "https://www.uniprot.org/uniprot/" + this.props.unpID + ".fasta";
    return (
      <Box>
        <Typography variant="h5" gutterBottom>
          3. Add variants
        </Typography>
        <Typography variant="overline" gutterBottom>
          <Link href={unpSeqURL} target="_blank" rel="noreferrer">
            canonical sequence
          </Link>{" "}
          length: {this.props.sequence.length}
        </Typography>
        {this.props.hasAnnotated && (
          <Box>
            <FormControlLabel
              control={<Checkbox onChange={this.handleChange} />}
              label="Include annotated variants"
            />
          </Box>
        )}

        <VariantInput
          variants={this.props.variants}
          sequence={this.props.sequence}
          setVariants={this.props.setVariants}
        />

        <ChipArray
          variants={this.props.variants}
          handleDelete={this.handleDelete}
        />
      </Box>
    );
  }
}
