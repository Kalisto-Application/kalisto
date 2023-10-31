import { useRef, useState } from 'react';
import { useBoolean, useOnClickOutside } from 'usehooks-ts';

import { models } from '../../wailsjs/go/models';
import { Menu, MenuItemProp, MenuProps } from './Menu';

import expandIcon from '../../assets/icons/expand.svg';

type fileListtProps = {
  items: itemProps[];
  onCloseInput: () => void;
  editFile: (value: string) => void;
};
type itemProps = {
  file: models.File;
  inEdit: boolean;
  isActive: boolean;
  onClick: () => void;
  menu: MenuItemProp[];
};
const FileList: React.FC<fileListtProps> = ({
  items,
  onCloseInput,
  editFile,
}) => {
  const [valueEdit, setValueEdit] = useState('');

  const onKeyDown = (e: React.KeyboardEvent) => {
    if (e.code == 'Enter' && valueEdit !== '') {
      editFile(valueEdit.trim());
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
          className={`relative flex cursor-pointer justify-between py-1 pl-10 pr-4 hover:bg-borderFill ${
            it.isActive ? 'bg-textBlockFill' : ''
          }`}
          onClick={it.onClick}
        >
          {it.inEdit ? (
            <input
              className="border-1 w-[75%] border-[1px] border-borderFill bg-textBlockFill px-3 placeholder:text-[14px] placeholder:text-secondaryText"
              type="text"
              onFocus={(e) => {
                e.target.select();
                setValueEdit(it.file.name);
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
          <SubMenu items={it.menu} />
        </li>
      ))}
    </ul>
  );
};
export default FileList;

const SubMenu: React.FC<MenuProps> = ({ items }) => {
  const subMenuRef = useRef(null);
  const { value, toggle, setFalse } = useBoolean(false);

  useOnClickOutside(subMenuRef, () => setFalse());

  return (
    <>
      {/* button submenu  */}
      <button onClick={toggle}>
        <img src={expandIcon} alt="" />
      </button>
      {/* Sub menu */}
      {value && (
        <div ref={subMenuRef} className="absolute right-2 top-9 w-[70%]">
          <Menu
            items={items.map((it) => {
              return {
                ...it,
                onClick: () => {
                  toggle();
                  it.onClick?.();
                },
              };
            })}
          />
        </div>
      )}
    </>
  );
};
