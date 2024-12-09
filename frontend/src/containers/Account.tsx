import { Button, Card, FormGroup, H3, H5, Icon, InputGroup } from "@blueprintjs/core";
import useAccountInfo, { useGetRole } from "../hooks/useAccountInfo.ts";
import React, { useEffect, useState } from "react";
import { LandingNavbar } from "../components/LandingNavbar.tsx";
import { InfoPage, PageTag, Role } from "../components/InfoPage.tsx";
import { Sidebar } from "../components/Sidebar.tsx";
import { useIsFounder } from "../hooks/useBusinessInfo.ts";
import { useNavigate } from "react-router-dom";
import { HelpDialog } from "../components/HelpDialog.tsx";

export const Account = () => {
    const [isOpen, setIsOpen] = useState(false);
    const navigate = useNavigate();
    const role = useGetRole();
    const handlePageChange = (page: PageTag) => {
        navigate("/" + page);
    };

    const handleDialog = (newOpen: boolean) => {
        setIsOpen(newOpen);
    }

    return (
        <div>
            <HelpDialog isOpen={isOpen} setIsOpen={setIsOpen} />
            <LandingNavbar handleDialog={handleDialog} showAccount={false} />
            <div className='App'>
                <Sidebar isAdminFounder={role === Role.Admin || role === Role.Founder} handlePageChange={handlePageChange} />
                <InfoPage />
            </div>
        </div>
    );

}

export default Account;