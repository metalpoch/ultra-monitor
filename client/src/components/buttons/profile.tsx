import { useEffect, useState } from 'react';
import type { User } from '../../models/user.ts';
import { Routes } from '../../constant/routes';

interface Props {
    user: User;
}

export default function ProfileButton({ user }: Props) {

    const [showDropdown, setShowDropdown] = useState(false);

    useEffect(() => {
        const dropdown = document.getElementById('dropdownAvatar');
        if (dropdown && showDropdown) {
            dropdown.classList.remove('hidden');
            dropdown.classList.add('visible');
            dropdown.classList.add('top-14');
        } else if (dropdown) {
            dropdown.classList.add('hidden');
            dropdown.classList.remove('visible');
            dropdown.classList.remove('top-14');
        }
    }, [showDropdown]);

    return (
        <div className="flex flex-col items-end gap-1 h-10 w-30 px-2">
            <button id="dropdownUserAvatarButton" data-dropdown-toggle="dropdownAvatar" className="w-8 h-8 flex text-sm rounded-full md:me-0 focus:ring-4 focus:ring-gray-300 dark:focus:ring-gray-600" type="button" onClick={() => setShowDropdown(!showDropdown)}>
                <img className="w-8 h-8 rounded-full" src="/assets/profile.svg" alt="user photo" />
            </button>
            <div id="dropdownAvatar" className="z-10 hidden absolute bg-white divide-y divide-gray-100 rounded-lg shadow w-44 dark:bg-gray-700 dark:divide-gray-600">
                <div className="px-4 py-3 text-sm text-gray-900 dark:text-white">
                    <div>{user.fullname}</div>
                    <div className="font-medium truncate">{user.email}</div>
                </div>
                <ul className="py-2 text-sm text-gray-700 dark:text-gray-200" aria-labelledby="dropdownUserAvatarButton">
                    <li>
                        <a href={Routes.PROFILE} className="block px-4 py-2 hover:bg-gray-100 dark:hover:bg-gray-600 dark:hover:text-white">Perfil</a>
                    </li>
                </ul>
                <div className="py-2">
                    <a href="#" className="block px-4 py-2 text-sm text-gray-700 hover:bg-gray-100 dark:hover:bg-gray-600 dark:text-gray-200 dark:hover:text-white">Cerrar Sesión</a>
                </div>
            </div>
        </div>
    );
}