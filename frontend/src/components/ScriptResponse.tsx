import React, { useContext } from "react";
import { Editor } from "../ui/Editor";
import { Context } from "../state";

export const ScriptResponse: React.FC = () => {
  const ctx = useContext(Context);

  return (
    <div className="bg-textBlockFill w-1/2">
      <Editor
        key={ctx.state.scriptResponse}
        value={ctx.state.scriptResponse}
        readonly
      />
    </div>
  );
};
