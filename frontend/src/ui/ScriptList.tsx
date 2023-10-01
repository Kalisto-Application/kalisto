import { useRef, useState } from 'react';
import { models } from '../../wailsjs/go/models';
import { Menu } from './Menu';
import { useOnClickOutside } from '../hooks/useOnClickOutside';
import plusIcon from '../../assets/icons/plus.svg';

import subMenuIcon from '../../assets/icons/subMenu.svg';

type scriptListProps = {
  addScript: (value: string) => void;
  activeWorkspace?: models.Workspace;
  setActiveScript: (id: string) => void;
  items: itemProps[];
};
interface itemProps {
  text: string;
  icon?: string;
  onClick?: (e: React.MouseEvent) => void;
}
const ScriptList: React.FC<scriptListProps> = ({
  addScript,
  activeWorkspace,
  items,
  setActiveScript: deleteScript,
}) => {
  return (
    <>
      <ScriptNewCreate
        addScript={addScript}
        activeWorkspace={activeWorkspace}
      />

      <ItemList
        deleteScript={(id) => deleteScript(id)}
        activeWorkspace={activeWorkspace}
        items={items}
      />
    </>
  );
};
export default ScriptList;

type ScriptNewProps = {
  addScript: (value: string) => void;
  activeWorkspace?: models.Workspace;
};

const ScriptNewCreate: React.FC<ScriptNewProps> = ({
  addScript,

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
    setIsMode(false);
  };
  return (
    <div className="mb-3 flex flex-col items-center ">
      {isMode ? (
        <button
          className="flex items-center gap-x-2  rounded-md border-[1px] border-borderFill bg-primaryFill px-3 py-1 transition duration-500 ease-in-out hover:bg-textBlockFill"
          onClick={onClick}
        >
          <img src={plusIcon} alt="" />
          <span className="text-lg">Add Script</span>
        </button>
      ) : (
        <>
          <div className="px-5">
            <input
              placeholder="Name script"
              className="border-1 mb-1 w-full border-[1px] border-borderFill  bg-textBlockFill px-3 placeholder:text-[14px] placeholder:text-secondaryText"
              type="text"
              autoFocus
              onChange={(e) => updateValue(e)}
              onKeyDown={onKeyDown}
            />
          </div>
          <div className="text-xs">Push Enter to rename </div>
        </>
      )}
      {valueValidate ? (
        <span className="text-[13px] text-red">
          A script name must not be empty
        </span>
      ) : null}
    </div>
  );
};

interface ItemListProps {
  activeWorkspace?: models.Workspace;
  deleteScript: (id: string) => void;
  items: itemProps[];
}

const ItemList: React.FC<ItemListProps> = ({
  activeWorkspace,
  deleteScript,
  items,
}) => {
  const [isMode, setIsMode] = useState(false);
  const [idSubMenu, setIdSubMenu] = useState('');

  const subMenuRef = useRef(null);

  useOnClickOutside(subMenuRef, () => setIsMode(false));

  const active = 'text-red';

  return (
    <ul className="text-center">
      {activeWorkspace?.scriptFiles ? (
        activeWorkspace?.scriptFiles
          .map((it, indx) => (
            <li
              key={indx}
              className=" text-ms relative  flex cursor-pointer justify-between px-3"
              onClick={() => {
                setIdSubMenu(it.id);
                deleteScript(it.id);
              }}
            >
              <button className={`${idSubMenu === it.id ? active : null}`}>
                {it.name}
              </button>
              {/* button submenu  */}
              <button
                onClick={(e) => {
                  setIsMode(true);
                }}
              >
                <img src={subMenuIcon} alt="" />
              </button>
              {/* Sub menu */}
              {isMode && idSubMenu === it.id ? (
                <div
                  ref={subMenuRef}
                  className="absolute right-2 top-5 w-[70%]"
                >
                  <Menu items={items} />
                </div>
              ) : null}
            </li>
          ))
          .reverse()
      ) : (
        <h2>No scripts found</h2>
      )}
    </ul>
  );
};
