import React, { useContext } from "react";
import { EditorSwitcher } from "../ui/EditorSwitcher";
import { CodeEditor } from "./CodeEditor";

import { Context } from "../state";

interface RequestEditorProps {
}

export const RequestEditor: React.FC<RequestEditorProps> = () => {
    const ctx = useContext(Context)
    const switchRequestEditor = (i: number) => void [
        ctx.dispatch({type: 'switchRequestEditor', i: i})
    ]

    const editors = [<CodeEditor key={0} text={ctx.state.requestText} type="changeRequestText"/>, 
    <CodeEditor key={1} text={ctx.state.requestMetaText} type="changeMetaText"/>]

    return (
        <div className="bg-textBlockFill w-1/2">
            <EditorSwitcher items={[
                {title: 'Request', onClick: switchRequestEditor},
                {title: 'Metadata', onClick: switchRequestEditor},
            ]} active={ctx.state.activeRequestEditor || 0}/>
            {editors[ctx.state.activeRequestEditor] || editors[0]}
        </div>
    );
}
