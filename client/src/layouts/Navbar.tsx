import { useEffect, useState } from "react";
import type { User } from '../models/user.ts';
import { Routes } from "../constant/routes";

interface Props {
    user: User;
}

export default function Navbar({ user }: Props) {

    const [showMenu, setShowMenu] = useState(false);
    const [showProfile, setShowProfile] = useState(false);
   
    useEffect(() => {
        const menuList = document.getElementById('menu-list');
        const dropdownButton = document.getElementById('dropdownUserAvatarButton');
        const dropdown = document.getElementById('dropdownAvatar');

        if (menuList && showMenu) {
            menuList.classList.remove('max-md:hidden');
            menuList.classList.add('max-md:visible');
        } else if (menuList) {
            menuList.classList.add('max-md:hidden');
            menuList.classList.remove('max-md:visible');
        }

        if (showProfile &&dropdownButton && dropdown) {
            dropdown.classList.remove('hidden');
            dropdown.classList.add('visible');
        } else if (dropdownButton && dropdown) {
            dropdown.classList.add('hidden');
            dropdown.classList.remove('visible');
        }
    }, [showMenu, showProfile]);

    return(
        <nav className="w-full h-fit py-4 flex flex-row justify-between items-center max-md:justify-end max-md:items-end max-md:px-4 max-md:gap-4">
            <div></div>
            <button id="menu" type="button" className="hidden max-md:block" onClick={() => setShowMenu(!showMenu)}>
                <img src="/assets/menu.svg" alt="menu" />
            </button>
            <ul id="menu-list" className="flex flex-row justify-evenly gap-4 max-md:w-full max-md:flex-col max-md:items-center max-md:hidden">
                <li className="w-32 px-1 py-1 flex flex-row justify-center font-light rounded-full transition-all duration-200 ease-in hover:text-gray-50 hover:bg-gray-300">
                    <a href={Routes.HOME}>Dashboard</a>
                </li>
                <li className="w-32 px-1 py-1 flex flex-row justify-center font-light rounded-full transition-all duration-200 ease-in hover:text-gray-50 hover:bg-gray-300">
                    <a href="#">OLT</a>
                </li>
                <li className="w-32 px-1 py-1 flex flex-row justify-center font-light rounded-full transition-all duration-200 ease-in hover:text-gray-50 hover:bg-gray-300">
                    <a href="#">Tendencia</a>
                </li>
                <li className="w-32 px-1 py-1 flex flex-row justify-center font-light rounded-full transition-all duration-200 ease-in hover:text-gray-50 hover:bg-gray-300">
                    <a href="#">RodolfIA</a>
                </li>
                <li className="w-32 px-1 py-1 flex flex-row justify-center font-light rounded-full transition-all duration-200 ease-in hover:text-gray-50 hover:bg-gray-300">
                    <a href="#">Reportes</a>
                </li>
            </ul>
            <div className="flex flex-col justify-center items-end gap-1 h-10 w-30 px-2">
                <button id="dropdownUserAvatarButton" data-dropdown-toggle="dropdownAvatar" className="w-8 h-8 flex text-sm rounded-full md:me-0 
                focus:ring-4 focus:ring-gray-300 dark:focus:ring-gray-600" type="button" onClick={() => setShowProfile(!showProfile)} onBlur={() => setShowProfile(false)}>
                    <img className="w-8 h-8 rounded-full" src="/assets/profile.svg" alt="user photo" />
                </button>
                <div id="dropdownAvatar" className="z-10 hidden absolute top-14 bg-white divide-y divide-gray-100 rounded-lg shadow w-44 dark:bg-gray-700 dark:divide-gray-600">
                    {user &&
                        <div className="px-4 py-3 text-sm text-gray-900 dark:text-white">
                            <div>{user.fullname}</div>
                            <div className="font-medium truncate">{user.email}</div>
                        </div>
                    }
                    <ul className="py-2 text-sm text-gray-700 dark:text-gray-200" aria-labelledby="dropdownUserAvatarButton">
                        {user &&
                            <li>
                                <a href={Routes.PROFILE} className="block px-4 py-2 hover:bg-gray-100 dark:hover:bg-gray-600 dark:hover:text-white">Perfil</a>
                            </li>
                        }
                    </ul>
                    <div className="py-2">
                        <a href="#" className="block px-4 py-2 text-sm text-gray-700 hover:bg-gray-100 dark:hover:bg-gray-600 dark:text-gray-200 dark:hover:text-white">Cerrar Sesión</a>
                    </div>
                </div>
            </div>            
        </nav>
    );
}