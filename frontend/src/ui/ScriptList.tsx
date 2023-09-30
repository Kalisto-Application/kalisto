import { useState } from 'react';
import { models } from '../../wailsjs/go/models';
import iconPlus from '../../assets/icons/plus.svg';
import editIcon from '../../assets/icons/edit.svg';
import deleteIcon from '../../assets/icons/delete.svg';

type scriptListProps = {
  setValidateWorkspace: (value: boolean) => void;
  addScript: (value: string) => void;
  activeWorkspace?: models.Workspace;
  deleteScript: (id: string) => void;
};

const ScriptList: React.FC<scriptListProps> = ({
  addScript,
  activeWorkspace,
  setValidateWorkspace,
  deleteScript,
}) => {
  return (
    <>
      <ScriptNewCreate
        addScript={addScript}
        activeWorkspace={activeWorkspace}
        setValidateWorkspace={(value) => setValidateWorkspace(value)}
      />

      <ul className="pl-10">
        <ItemList
          deleteScript={deleteScript}
          activeWorkspace={activeWorkspace}
        />
      </ul>
    </>
  );
};
export default ScriptList;

type ScriptNewProps = {
  setValidateWorkspace: (value: boolean) => void;
  addScript: (value: string) => void;
  activeWorkspace?: models.Workspace;
};

const ScriptNewCreate: React.FC<ScriptNewProps> = ({
  addScript,
  setValidateWorkspace,
  activeWorkspace,
}) => {
  const [value, setValue] = useState('');
  const [isMode, setIsMode] = useState(true);
  const [valueValidate, setValueValidate] = useState(false);

  const updateValue = (e: React.ChangeEvent<HTMLInputElement>) => {
    const newValue = e.target.value;
    setValueValidate(false);
    setValue(newValue);
  };

  const onKeyDown = (e: React.KeyboardEvent) => {
    if (e.code == 'Enter' && value === '') {
      setValueValidate(true);
      return;
    }

    if (e.code == 'Enter') {
      addScript(value);
      setIsMode(true);
      setValue('');
    }
  };

  const onClick = () => {
    if (!activeWorkspace) {
      setValidateWorkspace(true);
      return;
    }
    setIsMode(false);
  };
  return (
    <>
      {isMode ? (
        <button className="flex items-center gap-x-2 pl-3" onClick={onClick}>
          <img src={iconPlus} alt="" /> <span>New Script</span>
        </button>
      ) : (
        <div className=" text-center">
          <input
            className="border-1 mb-1 w-10/12 border-borderFill bg-textBlockFill"
            type="text"
            autoFocus
            onChange={(e) => updateValue(e)}
            onBlur={() => {
              setIsMode(true), setValueValidate(false);
            }}
            onKeyDown={onKeyDown}
          />
          <div className=" text-xs">Нажми Entre для создания</div>
        </div>
      )}
      {valueValidate ? (
        <span className="pl-3 text-center text-[13px] text-red">
          A script name must not be empty
        </span>
      ) : null}
    </>
  );
};

interface ItemListProps {
  activeWorkspace?: models.Workspace;
  deleteScript: (id: string) => void;
}

const ItemList: React.FC<ItemListProps> = ({
  activeWorkspace,
  deleteScript,
}) => {
  const items = [
    {
      // icon: editIcon,
      text: 'Edit',
      onClick: (e: React.MouseEvent) => {
        e.preventDefault();
      },
    },

    {
      // icon: deleteIcon,
      text: 'Delete',
      onClick: (e: React.SyntheticEvent) => {
        e.preventDefault();
      },
    },
  ];

  return (
    <>
      {activeWorkspace?.scriptFiles ? (
        activeWorkspace?.scriptFiles
          .map((it, indx) => (
            <li
              key={indx}
              className="text-ms hover: flex  cursor-pointer justify-between pr-3"
            >
              <button>{it.name}</button>
            </li>
          ))
          .reverse()
      ) : (
        <h2>Список пуст</h2>
      )}
    </>
  );
};
