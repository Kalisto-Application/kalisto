import { Menu, Transition } from '@headlessui/react';
import React, { ChangeEvent, Fragment, useRef, useState } from 'react';
import { Menu as M, MenuItem } from './Menu';

export type DropdownProps = {
  main: DropdownItemProps;
  items: DropdownItemProps[];
};

export type DropdownItemProps = {
  text: string;
  icon?: string;
  tip?: React.ReactElement;
  onClick?: (e: React.MouseEvent) => void;
  edit?: editable;
  menu?: DropdownItemProps[];
  divide?: boolean;
};

interface editable {
  inEdit: boolean;
  onEditDone: (newName: string) => void;
}

const shortifyMainText = (value: string): string => {
  return value.length > 15 ? value.slice(0, 12) + ' ...' : value;
};

export const Dropdown: React.FC<DropdownProps> = ({ main, items }) => {
  const mainText = shortifyMainText(main.text);

  return (
    <Menu>
      <Menu.Button className="flex h-[48px] shrink-0 flex-wrap content-center items-center gap-3 border-[1px] border-borderFill px-4 py-2.5">
        {main.icon && <img className="" src={main.icon} />}
        <p className="">{mainText}</p>
        {main.tip && <div className="ml-auto">{main.tip}</div>}
      </Menu.Button>
      <div className="relative mt-2 flex">
        <Transition
          as={Fragment}
          enter="transition ease-out duration-300"
          enterFrom="transform opacity-0 scale-95"
          enterTo="transform opacity-100 scale-100"
          leave="transition ease-in duration-75"
          leaveFrom="transform opacity-100 scale-100"
          leaveTo="transform opacity-0 scale-95"
        >
          <Menu.Items className="absolute z-10 flex cursor-pointer flex-col content-normal items-center justify-center rounded-md border-[1px] border-borderFill bg-primaryFill px-0 py-1">
            {items.map((it, i) => (
              <DropdownItem key={i} {...it} />
            ))}
          </Menu.Items>
        </Transition>
      </div>
    </Menu>
  );
};

const DropdownItem: React.FC<DropdownItemProps> = ({
  text,
  icon,
  tip,
  onClick,
  edit,
  menu,
  divide,
}) => {
  return (
    <Menu.Item as={Fragment}>
      {(render) => {
        if (edit?.inEdit) {
          return <EditableItem value={text} onDone={edit?.onEditDone} />;
        }
        let el = (
          <div
            onClick={onClick}
            className="ui-active:text-white relative flex h-11 w-[259.5px] flex-row items-center gap-x-[42px] px-4 py-2.5  ui-active:bg-textBlockFill"
          >
            {icon && <img src={icon} />}
            <a>{text}</a>
            {tip}
            {render.active && (
              <div className="absolute left-[200px] top-0 z-20 w-1/2">
                {menu && (
                  <M
                    items={menu.map((menu) => ({
                      text: menu.text,
                      icon: menu.icon,
                      onClick: menu.onClick,
                    }))}
                  />
                )}
              </div>
            )}
          </div>
        );
        return divide ? <div className="divide">{el}</div> : el;
      }}
    </Menu.Item>
  );
};

export default Dropdown;

interface EditableItemProps {
  value: string;
  onDone: (value: string) => void;
}

const EditableItem: React.FC<EditableItemProps> = ({ value, onDone }) => {
  const [editing, setEditing] = useState(value);

  const onKeyDown = (e: React.KeyboardEvent) => {
    if (e.code == 'Space') {
      e.preventDefault();
      setEditing((prev) => prev + ' ');
    }
    if (e.code == 'Enter') {
      onDone(editing);
    }
  };

  return (
    <div
      onClick={(e) => e.preventDefault()}
      className=" h-11 w-[259.5px] bg-textBlockFill leading-6"
    >
      <input
        className="    h-full w-full bg-textBlockFill text-secondaryText"
        value={editing}
        onChange={(e) => setEditing(e.target.value)}
        autoFocus={true}
        onKeyDown={onKeyDown}
        onBlur={(e) => {
          if (e.currentTarget === e.target) {
            onDone(editing);
          }
        }}
      />
      <div className=" bg-primaryFill pl-1 pt-1 text-xs">
        Push Enter to rename
      </div>
    </div>
  );
};
