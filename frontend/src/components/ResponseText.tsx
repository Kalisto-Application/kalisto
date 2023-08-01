import React, { useContext } from "react";
import { EditorSwitcher } from "../ui/EditorSwitcher";
import { Editor } from "../ui/Editor";
import { Context } from "../state";

interface props {
  body: string;
  meta: string;
}

export const ResponseText: React.FC<props> = ({body, meta}) => {
    const ctx = useContext(Context)
    const switchResponseEditor = (i: number) => void [
        ctx.dispatch({type: 'switchResponseEditor', i: i})
    ]
    const bodyKey = `request:${ctx.state.activeMethod?.fullName}:${body}`
    const metaKey = `meta:${ctx.state.activeMethod?.fullName}:${meta}`

    const editors = [<Editor key={bodyKey} value={body} readonly />, 
    <Editor key={metaKey} value={meta} readonly />]

    return (
        <div className="bg-textBlockFill w-1/2">
            <EditorSwitcher items={[
                {title: 'Response', onClick: switchResponseEditor},
                {title: 'Metadata', onClick: switchResponseEditor},
            ]} active={ctx.state.activeResponseEditor || 0}/>
            {editors[ctx.state.activeResponseEditor] || editors[0]}
        </div>
    );
}
