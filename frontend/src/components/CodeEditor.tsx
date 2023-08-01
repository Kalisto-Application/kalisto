import React, {ChangeEvent, useContext, useMemo } from "react";
import { Context } from "../state";
import { debounce, Action } from "../pkg";

import { Editor } from "../ui/Editor";

interface CodeEditorProps {
    text: string;
    type: 'changeRequestText' | 'changeMetaText' | 'changeVariables';
    action?: Action;
};

export const CodeEditor: React.FC<CodeEditorProps> = ({ text, type, action }) => {
    const ctx = useContext(Context);

    let debouncedAction: Action | undefined;
    if (action) {
      debouncedAction = useMemo<Action>(()=> {
        return debounce(action, 400);
      }, [])
    }

    let onChange = (value: string) => {
        ctx.dispatch({type: type, text: value});
        if (debouncedAction) {
          debouncedAction(value)
        }
    }

    return (
        <div>
          <Editor key={ctx.state.activeMethod?.fullName || ""} value={text} onChange={onChange}/>
        </div>
      );
}
