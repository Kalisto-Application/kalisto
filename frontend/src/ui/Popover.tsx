import { Popover, Transition } from '@headlessui/react';

type popoverUiType = {
  children: React.ReactNode;
  nameButton: string;
  text: string;
};

const PopoverUI: React.FC<popoverUiType> = ({ children, nameButton, text }) => {
  return (
    <Popover className="relative">
      {({ open, close }) => (
        <>
          <Popover.Button
            as="div"
            onClick={() => setTimeout(() => close(), 900)}
            className="flex"
          >
            {children}
          </Popover.Button>

          <Transition
            enter="transition duration-100 ease-out"
            enterFrom="transform scale-95 opacity-0"
            enterTo="transform scale-100 opacity-100"
            leave="transition duration-75 ease-out"
            leaveFrom="transform scale-100 opacity-100"
            leaveTo="transform scale-95 opacity-0"
          >
            <Popover.Panel className="  absolute ">
              <div className="rounded bg-blind p-1">{text}</div>
            </Popover.Panel>
          </Transition>
        </>
      )}
    </Popover>
  );
};

export default PopoverUI;
