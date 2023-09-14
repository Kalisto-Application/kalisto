import React, { useContext } from 'react';
import { EditorSwitcher } from '../ui/EditorSwitcher';
import { Editor } from '../ui/Editor';
import { Context } from '../state';

export const ResponseText: React.FC = () => {
  const ctx = useContext(Context);

  const switchResponseEditor = (i: number) =>
    void [ctx.dispatch({ type: 'switchResponseEditor', i: i })];
  const bodyKey = `request:${ctx.state.activeMethod?.fullName}:${ctx.state.response?.body}`;
  const metaKey = `meta:${ctx.state.activeMethod?.fullName}:${ctx.state.response?.metaData}`;

  const editors = [
    <Editor key={bodyKey} value={ctx.state.response?.body || ''} readonly />,
    <Editor
      key={metaKey}
      value={ctx.state.response?.metaData || ''}
      readonly
    />,
  ];

  return (
    <div className="w-1/2 bg-textBlockFill">
      <EditorSwitcher
        items={[
          { title: 'Response', onClick: switchResponseEditor },
          { title: 'Metadata', onClick: switchResponseEditor },
        ]}
        active={ctx.state.activeResponseEditor || 0}
      />
      {editors[ctx.state.activeResponseEditor] || editors[0]}
    </div>
  );
};
