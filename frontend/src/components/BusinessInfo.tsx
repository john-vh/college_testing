import { Button, Card, FormGroup, H3, H5, Icon, InputGroup, Intent, Tag, TextArea } from "@blueprintjs/core";
import useAccountInfo, { useGetRole } from "../hooks/useAccountInfo.ts";
import React, { useEffect, useMemo, useState } from "react";
import { LandingNavbar } from "../components/LandingNavbar.tsx";
import { BusinessInfo, useBusinessInfo, useUpdateBusiness } from "../hooks/useBusinessInfo.ts";
import { AddBusiness } from "./AddBusiness.tsx";
import { Role } from "./InfoPage.tsx";
import { useApproveBusiness } from "../hooks/useApproveBusiness.ts";

export const BusinessInfoPage = () => {
    const isAdmin = useGetRole() === Role.Admin;
    const { businessInfo, fetchData } = useBusinessInfo({ isAdmin });
    const approveBusiness = useApproveBusiness();
    const { updateBusiness } = useUpdateBusiness();
    const [businessAdd, setBusinessAdd] = useState(false);
    const [isReadonly, setIsReadonly] = useState<boolean[]>([]);
    const [newData, setNewData] = useState<BusinessInfo[]>([]);

    const delay = (ms: number) => new Promise(resolve => setTimeout(resolve, ms));

    const addNewBusiness = () => {
        setBusinessAdd(true); // Change state to show the Add Business page
    };

    const handleEditBusiness = (index: number) => {
        const newIsReadonly = [...isReadonly];
        newIsReadonly[index] = false;
        console.log(newIsReadonly);
        setIsReadonly(newIsReadonly);
    }

    const handleInputChange = (index: number, key: string, value: string) => {
        const newData = [...businessInfo];
        newData[index][key] = value;
        setNewData(newData);
    }

    const handleSaveBusiness = (entry: BusinessInfo) => {
        updateBusiness(entry);
        console.log("updated");
    }

    const handleApproveBusiness = async (id: string) => {
        approveBusiness(id);
        await delay(100);
        await fetchData();
        // console.log("approved");
    }

    const sortedBusinessInfo = useMemo(
        () => [...businessInfo].sort((a, b) => a.name.localeCompare(b.name)),
        [businessInfo]
    );

    useEffect(() => {
        if (businessInfo && businessInfo.length > 0) {
            setIsReadonly(Array(businessInfo.length).fill(true));
            setNewData(Array(businessInfo.length));
        }
    }, [businessInfo]);

    if (isAdmin) {
        businessInfo.sort((a, b) => {
            if (a.status === b.status) {
                return 0;
            }
            return a.status === "pending" ? -1 : 1;
        });
    }

    if (businessAdd) {
        return <AddBusiness setBusinessAdd={setBusinessAdd} />
    }

    return (
        <div className="content">
            <Button
                intent="primary"
                large={true}
                style={{ width: "100%", marginBottom: "20px" }}
                onClick={() => addNewBusiness()}
            >
                Add Business
            </Button>
            <div style={{ display: 'flex', flexDirection: 'column', gap: '20px' }}>
                {sortedBusinessInfo != null && sortedBusinessInfo.map((entry, index) => (
                    <Card interactive={true} >
                        <div className="Flex" style={{ justifyContent: "space-between" }}>
                            <H3>Business Information</H3>
                            <Tag round intent={entry.status === "active" ? Intent.SUCCESS : Intent.WARNING}>{entry.status.toLocaleUpperCase()}</Tag>
                        </div>

                        <FormGroup label="Name"
                            labelFor="name" >
                            <InputGroup asyncControl id="name" onChange={(e) => handleInputChange(index, "name", e.target.value)} value={entry.name} readOnly={isReadonly[index]} />
                        </FormGroup>
                        <FormGroup label="Website"
                            labelFor="website" >
                            <InputGroup asyncControl id="website" onChange={(e) => handleInputChange(index, "website", e.target.value)} value={entry.website} readOnly={isReadonly[index]} />
                        </FormGroup>
                        <FormGroup label="Description"
                            labelFor="desc" >
                            <TextArea asyncControl id="desc" onChange={(e) => handleInputChange(index, "desc", e.target.value)} value={entry.desc} readOnly={isReadonly[index]} fill />
                        </FormGroup>
                        {(isAdmin && entry.status === "pending") && <Button intent="success" style={{ marginRight: "10px" }} onClick={() => handleApproveBusiness(entry.id)}>Approve business</Button>}
                        {isReadonly[index] ? <Button onClick={() => handleEditBusiness(index)}>Manage business</Button> : <Button onClick={() => handleSaveBusiness(entry)}>Save changes</Button>}
                    </Card>
                ))}
            </div>
        </div>
    );
}
