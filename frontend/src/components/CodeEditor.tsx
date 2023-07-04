import React, {ChangeEvent, useContext, useEffect} from "react";
import { Context } from "../state";

interface CodeEditorProps {
    text: string;
    type: 'changeRequestText' | 'changeMetaText';
}

export const CodeEditor: React.FC<CodeEditorProps> = ({ text, type }) => {
    const ctx = useContext(Context);

    const onChange = (e: ChangeEvent<HTMLTextAreaElement>) => {
        ctx.dispatch({type: type, text: e.target.value});
    }

    return (
        <div>
          <textarea value={text} onChange={onChange} className="w-[480px] h-[600px] bg-codeSectionBg text-inputPrimary"/>
        </div>
      );
}
