import { useRef, useState } from 'react';
import { models } from '../../wailsjs/go/models';
import { Menu } from './Menu';
import { useOnClickOutside } from 'usehooks-ts';
import plusIcon from '../../assets/icons/plus.svg';

import subMenuIcon from '../../assets/icons/subMenu.svg';

type fileListtProps = {
  addFile: (value: string) => void;
  activeWorkspace?: models.Workspace;
  setActiveScript: (id: string) => void;
  items: itemProps[];
  isOpenEditInput: boolean;
  onCloseInput: () => void;
  editFile: (value: string) => void;
  isModeSubMenu: boolean;
  closeSubMenu: () => void;
  openSubMenu: () => void;
  activeScript: string;
};
interface itemProps {
  text: string;
  icon?: string;
  onClick?: (e: React.MouseEvent) => void;
}
const FileList: React.FC<fileListtProps> = ({
  addFile,
  activeWorkspace,
  items,
  setActiveScript,
  isOpenEditInput,
  onCloseInput,
  editFile,
  isModeSubMenu,
  closeSubMenu,
  openSubMenu,
  activeScript,
}) => {
  return (
    <>
      <ScriptNewCreate addFile={addFile} />

      <ItemList
        setActiveScript={(id) => setActiveScript(id)}
        activeWorkspace={activeWorkspace}
        items={items}
        isOpenEditInput={isOpenEditInput}
        onCloseInput={onCloseInput}
        editFile={editFile}
        isModeSubMenu={isModeSubMenu}
        closeSubMenu={closeSubMenu}
        openSubMenu={openSubMenu}
        activeScript={activeScript}
      />
    </>
  );
};
export default FileList;

type ScriptNewProps = {
  addFile: (value: string) => void;
};

const ScriptNewCreate: React.FC<ScriptNewProps> = ({ addFile }) => {
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
      addFile(value);
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

type ItemListProps = {
  activeWorkspace?: models.Workspace;
  setActiveScript: (id: string) => void;
  items: itemProps[];
  isOpenEditInput: boolean;
  onCloseInput: () => void;
  editFile: (value: string) => void;
  isModeSubMenu: boolean;
  closeSubMenu: () => void;
  openSubMenu: () => void;
  activeScript: string;
};

const ItemList: React.FC<ItemListProps> = ({
  activeWorkspace,
  setActiveScript,
  items,
  isOpenEditInput,
  onCloseInput,
  editFile,
  isModeSubMenu,
  closeSubMenu,
  openSubMenu,
  activeScript,
}) => {
  const [idSubMenu, setIdSubMenu] = useState(activeScript);
  const [valueEdit, setValueEdit] = useState('');
  const subMenuRef = useRef(null);

  useOnClickOutside(subMenuRef, () => closeSubMenu());

  const active = 'text-red';

  const onKeyDown = (e: React.KeyboardEvent) => {
    if (e.code == 'Enter' && valueEdit !== '') {
      editFile(valueEdit);
      onCloseInput();
      setValueEdit('');
    }
  };

  const updateValueEdit = (e: React.ChangeEvent<HTMLInputElement>) => {
    setValueEdit(e.target.value);
  };
  return (
    <ul className="text-center">
      {activeWorkspace?.scriptFiles.length !== 0 ? (
        activeWorkspace?.scriptFiles.map((it, indx) => (
          <li
            key={indx}
            className=" text-ms relative  flex cursor-pointer justify-between px-3"
            onClick={() => {
              setIdSubMenu(it.id);
              setActiveScript(it.id);
            }}
          >
            {isOpenEditInput && idSubMenu === it.id ? (
              <input
                className="border-1   w-[75%] border-[1px]  border-borderFill bg-textBlockFill px-3 placeholder:text-[14px] placeholder:text-secondaryText"
                type="text"
                onFocus={(e) => {
                  e.target.select();
                }}
                autoFocus
                value={valueEdit || it.name}
                onChange={updateValueEdit}
                onBlur={onCloseInput}
                onKeyDown={onKeyDown}
              />
            ) : (
              <button className={`${idSubMenu === it.id ? active : ''}`}>
                {it.name}
              </button>
            )}
            {/* button submenu  */}
            <button onClick={openSubMenu}>
              <img src={subMenuIcon} alt="" />
            </button>
            {/* Sub menu */}
            {isModeSubMenu && idSubMenu === it.id ? (
              <div ref={subMenuRef} className="absolute right-2 top-5 w-[70%]">
                <Menu items={items} />
              </div>
            ) : null}
          </li>
        ))
      ) : (
        <h2>No scripts found</h2>
      )}
    </ul>
  );
};
