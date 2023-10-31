import { useRef, useState } from 'react';
import { useOnClickOutside } from 'usehooks-ts';

import { models } from '../../wailsjs/go/models';
import { Menu } from './Menu';

import expandIcon from '../../assets/icons/expand.svg';

type fileListtProps = {
  setActiveScript: (id: string) => void;
  items: itemProps[];
  isOpenEditInput: boolean;
  onCloseInput: () => void;
  editFile: (value: string) => void;
  isModeSubMenu: string;
  activeScriptId: string;
  setIsModeSubMenu: (value: string) => void;
};
type itemProps = {
  file: models.File;
  menu: menuProps[];
};
type menuProps = {
  text: string;
  icon?: string;
  onClick?: (e: React.MouseEvent) => void;
};

const FileList: React.FC<fileListtProps> = ({
  setActiveScript,
  items,
  isOpenEditInput,
  onCloseInput,
  editFile,
  isModeSubMenu,
  activeScriptId,
  setIsModeSubMenu,
}) => {
  const [isActive, setIsActive] = useState(activeScriptId);
  const [valueEdit, setValueEdit] = useState('');

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
    <ul className="flex-1">
      {items.map((it, indx) => (
        <li
          key={indx}
          className={`relative flex h-[32px] cursor-pointer justify-between pl-10 pr-4 hover:bg-borderFill ${
            it.file.id === isActive ? 'bg-textBlockFill' : ''
          }`}
          onClick={() => {
            setIsActive(it.file.id);
          }}
        >
          {isOpenEditInput && isActive === it.file.id ? (
            <input
              className="border-1 w-[75%] border-[1px] border-borderFill bg-textBlockFill px-3 placeholder:text-[14px] placeholder:text-secondaryText"
              type="text"
              onFocus={(e) => {
                e.target.select();
              }}
              autoFocus
              value={valueEdit || it.file.name}
              onChange={updateValueEdit}
              onBlur={onCloseInput}
              onKeyDown={onKeyDown}
            />
          ) : (
            <span className="text-right font-[Inter]">{it.file.name}</span>
          )}
          <SubMenu
            setActiveScript={setActiveScript}
            isModeSubMenu={isModeSubMenu}
            activeScriptId={activeScriptId}
            item={it}
            setIsModeSubMenu={setIsModeSubMenu}
          />
        </li>
      ))}
    </ul>
  );
};
export default FileList;

type propsSubMenu = {
  setActiveScript: (id: string) => void;
  isModeSubMenu: string;
  activeScriptId: string;
  item: {
    file: models.File;
    menu: menuProps[];
  };
  setIsModeSubMenu: (value: string) => void;
};

const SubMenu: React.FC<propsSubMenu> = ({
  setActiveScript,
  activeScriptId,
  isModeSubMenu,
  item,
  setIsModeSubMenu,
}) => {
  const subMenuRef = useRef(null);

  useOnClickOutside(subMenuRef, () => setIsModeSubMenu(''));

  return (
    <>
      {/* button submenu  */}
      <button
        onClick={(e) => {
          setIsModeSubMenu(item.file.id);
          setActiveScript(item.file.id);
          e.stopPropagation();
        }}
      >
        <img src={expandIcon} alt="" />
      </button>
      {/* Sub menu */}
      {isModeSubMenu && activeScriptId === item.file.id ? (
        <div ref={subMenuRef} className="absolute right-2 top-9 w-[70%]">
          <Menu items={item.menu} />
        </div>
      ) : null}
    </>
  );
};
