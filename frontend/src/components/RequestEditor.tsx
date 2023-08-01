import React, { useContext } from "react";
import { EditorSwitcher } from "./EditorSwitcher";
import { CodeEditor } from "./CodeEditor";

import { Context } from "../state";

interface RequestEditorProps {
}

export const RequestEditor: React.FC<RequestEditorProps> = () => {
    const ctx = useContext(Context)

    const editors = [<CodeEditor key={0} text={ctx.state.requestText} type="changeRequestText"/>, 
    <CodeEditor key={1} text={ctx.state.requestMetaText} type="changeMetaText"/>]

    return (
        <div className="bg-textBlockFill">
            <EditorSwitcher switches={['Request', 'Metadata']}/>
            {editors[ctx.state.activeEditor] || editors[0]}
        </div>
    );
}
