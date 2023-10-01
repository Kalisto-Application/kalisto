import React, { useContext, useState } from 'react';
import { Context } from '../state';
import { models } from '../../wailsjs/go/models';
import DeletePopup from './DeletePopup';
import FileList from '../ui/FileList';
import { UpdateWorkspace } from '../../wailsjs/go/api/Api';
import editIcon from '../../assets/icons/edit.svg';
import deleteIcon from '../../assets/icons/delete.svg';

const ScriptCollectionView: React.FC = () => {
  const ctx = useContext(Context);
  const [activeScript, setActiveScript] = useState('');
  const [isOpenDeletePopup, setIsOpenDeletePopup] = useState('');

  let textScript = ctx.state.scriptText;

  console.log('textScript', textScript);

  const addFile = (value: string) => {
    let workspace = new models.Workspace({
      ...ctx.state.activeWorkspace,
      scriptFiles: [
        ...(ctx.state.activeWorkspace?.scriptFiles || []),
        {
          name: value,
          createdAt: new Date(),
          content: '',
          headers: '',
          id: value,
        },
      ],
    });

    UpdateWorkspace(workspace).then((res) =>
      ctx.dispatch({
        type: 'updateWorkspace',
        workspace: workspace,
      })
    );
  };
  const deleteFile = (id: string) => {
    let workspace = new models.Workspace({
      ...ctx.state.activeWorkspace,
      scriptFiles: ctx.state.activeWorkspace?.scriptFiles.filter(
        (s) => s.id !== activeScript
      ),
    });
    UpdateWorkspace(workspace).then((res) => {
      ctx.dispatch({
        type: 'updateWorkspace',
        workspace: workspace,
      });
    });
  };

  const items = [
    {
      icon: editIcon,
      text: 'Edit',
      onClick: () => {},
    },

    {
      icon: deleteIcon,
      text: 'Delete',
      onClick: () => {
        setIsOpenDeletePopup(activeScript);
      },
    },
  ];

  return (
    <>
      {ctx.state.activeWorkspace ? (
        <>
          <FileList
            addFile={(value: string) => addFile(value)}
            activeWorkspace={ctx.state.activeWorkspace}
            setActiveScript={(id: string) => setActiveScript(id)}
            items={items}
          />
          <DeletePopup
            id={isOpenDeletePopup}
            isOpen={isOpenDeletePopup !== ''}
            onClose={() => setIsOpenDeletePopup('')}
            deleteScript={() => deleteFile(activeScript)}
            title="Delete script?"
          />
        </>
      ) : null}
    </>
  );
};

export default ScriptCollectionView;
