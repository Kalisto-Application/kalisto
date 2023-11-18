import { useRef, useState } from 'react';
import { useBoolean, useOnClickOutside } from 'usehooks-ts';

import { models } from '../../wailsjs/go/models';
import { Menu, MenuItemProp, MenuProps } from './Menu';

import expandIcon from '../../assets/icons/expand.svg';

type fileListtProps = {
  items: itemProps[];
  onCloseInput: () => void;
  editFile: (value: string) => void;
  gIcon?: string;
};
type itemProps = {
  file: models.File;
  inEdit: boolean;
  isActive?: boolean;
  onClick?: () => void;
  menu: MenuItemProp[];
};
const FileList: React.FC<fileListtProps> = ({
  items,
  onCloseInput,
  editFile,
  gIcon,
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
    <ul className="">
      {items.map((it, indx) => (
        <li
          key={indx}
          className={`relative ${
            gIcon ? '' : 'mx-2'
          } flex     cursor-pointer px-4 py-1 hover:bg-borderFill  hover:rounded-3xl${
            it.isActive ? ' rounded-3xl bg-textBlockFill ' : ''
          }`}
          onClick={it.onClick}
        >
          {it.inEdit ? (
            <input
              className="border-1 w-[100%] border-[1px] border-borderFill bg-textBlockFill px-3 placeholder:text-[14px] placeholder:text-secondaryText"
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
            <>
              {' '}
              {gIcon && <img src={gIcon} className="mr-2.5" />}{' '}
              <div className=" w-full font-[Inter]">{it.file.name}</div>
            </>
          )}
          <SubMenu items={it.menu} />
        </li>
      ))}
    </ul>
  );
};
export default FileList;

const SubMenu: React.FC<MenuProps> = ({ items }) => {
  const { value, toggle, setFalse } = useBoolean(false);
  const subMenuRef = useRef(null);

  useOnClickOutside(subMenuRef, () => setFalse());

  return (
    <div className="ml-1 flex w-[22px] justify-end" ref={subMenuRef}>
      {/* button submenu  */}
      <button
        onClick={(e) => {
          toggle();
          e.stopPropagation();
        }}
      >
        <img src={expandIcon} alt="" />
      </button>
      {/* Sub menu */}
      {value && (
        <div className="absolute right-2 top-9  ">
          <Menu
            items={items.map((it) => {
              return {
                ...it,
                onClick: (e: React.MouseEvent) => {
                  e.stopPropagation();
                  toggle();
                  it.onClick?.(e);
                },
              };
            })}
          />
        </div>
      )}
    </div>
  );
};
