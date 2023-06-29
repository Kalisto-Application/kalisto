import React, {ChangeEvent, useEffect} from "react";

interface CodeEditorProps {
    text: string;
    setText: (t :string) => void;
}

export const CodeEditor: React.FC<CodeEditorProps> = ({text, setText}) => {
    const onChange = (e: ChangeEvent<HTMLTextAreaElement>) => {
        setText(e.target.value)
    }
    console.log('text: ', text)

    return (
        <div>
          <textarea value={text} onChange={onChange} className="w-[480px] h-[600px] bg-codeSectionBg text-inputPrimary"/>
        </div>
      );
}
