import classNames from "classnames";
import { useRouter } from "next/router";
import { useState } from "react";
import { useAuth } from "../../store/auth-context";
import Backdrop from "../backdrop/backdrop";
import LoginModal from "../login-modal/login-modal";

function Header() {
  const router = useRouter();
  const { isAuth, data } = useAuth();
  const [profileModalIsOpen, setProfileModalIsOpen] = useState(false);

  const handleClickProfile = () => {
    setProfileModalIsOpen(true);
  };

  const handleClickBackdrop = () => {
    setProfileModalIsOpen(false);
  }

  return (
    <>
      <nav className="bg-white sticky top-0 z-30 border-gray-500 border-b px-2 sm:px-4 py-2.5 dark:bg-gray-800">
        <div className="container flex flex-wrap justify-between items-center mx-auto">
          <a href="/" className="flex">
            <span className="self-center text-xl font-semibold whitespace-nowrap dark:text-white">
              Patrick Arvatu
            </span>
          </a>
          <div className="flex items-center md:order-2">
            <button
              onClick={handleClickProfile}
              type="button"
              className="flex mr-3 text-sm bg-gray-800 rounded-full md:mr-0 focus:ring-4 focus:ring-gray-300 dark:focus:ring-gray-600"
              aria-expanded="false"
            >
              <img
                className="w-8 h-8 rounded-full"
                src={isAuth ? data.pfpUrl : "/image/blank-user.png"}
                alt="user photo"
              />
            </button>
          </div>
          <div
            className="hidden justify-between items-center w-full md:flex md:w-auto md:order-1"
            id="mobile-menu-2"
          >
            <ul className="flex flex-col mt-4 md:flex-row md:space-x-8 md:mt-0 md:text-sm md:font-medium">
              <li>
                <a
                  href="/"
                  className={classNames({
                    block: true,
                    "py-2": true,
                    "pr-4": true,
                    "pl-3": true,
                    "text-white": true,
                    "bg-blue-700": true,
                    rounded: true,
                    "md:bg-transparent": true,
                    "md:text-blue-700": router.pathname === "/",
                    "md:p-0": true,
                    "dark:text-white": true,
                  })}
                  aria-current="page"
                >
                  Home
                </a>
              </li>
              <li>
                <a
                  href="/about"
                  className={classNames({
                    block: true,
                    "py-2": true,
                    "pr-4": true,
                    "pl-3": true,
                    "text-white": true,
                    "bg-blue-700": true,
                    rounded: true,
                    "md:bg-transparent": true,
                    "md:text-blue-700": router.pathname === "/about",
                    "md:p-0": true,
                    "dark:text-white": true,
                  })}
                >
                  About
                </a>
              </li>
              <li>
                <a
                  href="/blog"
                  className={classNames({
                    block: true,
                    "py-2": true,
                    "pr-4": true,
                    "pl-3": true,
                    "text-white": true,
                    "bg-blue-700": true,
                    rounded: true,
                    "md:bg-transparent": true,
                    "md:text-blue-700": router.pathname === "/blog",
                    "md:p-0": true,
                    "dark:text-white": true,
                  })}
                >
                  Blog
                </a>
              </li>
              <li>
                <a
                  href="/projects"
                  className={classNames({
                    block: true,
                    "py-2": true,
                    "pr-4": true,
                    "pl-3": true,
                    "text-white": true,
                    "bg-blue-700": true,
                    rounded: true,
                    "md:bg-transparent": true,
                    "md:text-blue-700": router.pathname === "/projects",
                    "md:p-0": true,
                    "dark:text-white": true,
                  })}
                >
                  Projects
                </a>
              </li>
            </ul>
          </div>
        </div>
      </nav>
      {profileModalIsOpen && <LoginModal />}
      {profileModalIsOpen && <Backdrop onClick={handleClickBackdrop} />}
    </>
  );
}

export default Header;
