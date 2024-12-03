import { Button, Card, FormGroup, H3, H5, Icon, InputGroup } from "@blueprintjs/core";
import useAccountInfo, { useGetRole } from "../hooks/useAccountInfo.ts";
import React, { useEffect } from "react";
import { LandingNavbar } from "../components/LandingNavbar.tsx";
import { InfoPage, PageTag } from "../components/InfoPage.tsx";
import { Sidebar } from "../components/Sidebar.tsx";
import { useIsFounder } from "../hooks/useBusinessInfo.ts";

export const Account = () => {
    const [currentPage, setCurrentPage] = React.useState<PageTag>(PageTag.Account);
    const role = useGetRole();
    const handlePageChange = (page: PageTag) => {
        setCurrentPage(page);
    };

    return (
        <div>
            <LandingNavbar showAccount={false} />
            <div className='App'>
                <Sidebar handlePageChange={handlePageChange} />
                <InfoPage page={currentPage} role={role} />
            </div>
        </div>
    );

}

export default Account;