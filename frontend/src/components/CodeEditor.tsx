import React, {ChangeEvent, useContext, useMemo } from "react";
import { Context } from "../state";
import { debounce, Action } from "../pkg";

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

    let onChange = (e: ChangeEvent<HTMLTextAreaElement>) => {
        ctx.dispatch({type: type, text: e.target.value});
        if (debouncedAction) {
          debouncedAction(e.target.value)
        }
    }

    return (
        <div>
          <textarea value={text} onChange={onChange} className="w-[480px] h-[600px] bg-codeSectionBg text-inputPrimary"/>
        </div>
      );
}
