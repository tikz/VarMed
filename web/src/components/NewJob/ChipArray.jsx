import { Chip } from "@material-ui/core";
import React from "react";
import "../../styles/components/chip-array.scss";

export default function ChipArray(props) {
  return (
    <div className="chip-array">
      {props.variants.map((data) => {
        return (
          <Chip
            key={data.key}
            label={data.label}
            onDelete={() => props.handleDelete(data)}
            size="small"
          />
        );
      })}
    </div>
  );
}
