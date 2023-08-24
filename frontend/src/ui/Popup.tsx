import { Dialog } from '@headlessui/react';
import close from '../../assets/icons/close.svg';
interface propsPopup {
  children?: React.ReactNode;
  isOpen: boolean;
  onClose: () => void;
  title: string;
}

const Popup: React.FC<propsPopup> = ({ children, isOpen, onClose, title }) => {
  return (
    <Dialog open={isOpen} onClose={onClose} className="relative z-50">
      <div className="fixed inset-0 bg-[#00000062]" aria-hidden="true" />
      <div className="fixed inset-0  flex items-center justify-center p-4">
        <Dialog.Panel className="flex-[0_1_49.063rem] rounded-md border-2 border-textBlockFill bg-primaryFill p-8">
          <div className=" mb-10 flex items-center justify-between ">
            <h1 className="  text-lg font-bold leading-6 tracking-[1px]">
              {title}
            </h1>
            <button onClick={onClose}>
              <img src={close} alt="close" />
              {/* <span className="iconClose"></span> */}
            </button>
          </div>

          {children}
        </Dialog.Panel>
      </div>
    </Dialog>
  );
};
export default Popup;
