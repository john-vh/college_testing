import { Button, Card, FormGroup, H3, H2, H5, Icon, InputGroup, Intent, Tag, TextArea, FileInput } from "@blueprintjs/core";
import useAccountInfo, { useGetRole } from "../hooks/useAccountInfo.ts";
import React, { useEffect, useMemo, useState } from "react";
import { LandingNavbar } from "../components/LandingNavbar.tsx";
import { BusinessInfo, useBusinessInfo, useUpdateBusiness } from "../hooks/useBusinessInfo.ts";
import { AddBusiness } from "./AddBusiness.tsx";
import { Role } from "./InfoPage.tsx";
import { useApproveBusiness } from "../hooks/useApproveBusiness.ts";
import { useUploadImage } from "../hooks/useUploadImage.ts";

interface BusinessInfoProps {
    isAdmin: boolean;
}

export const BusinessInfoPage = ({ isAdmin = false }: BusinessInfoProps) => {
    const { businessInfo, fetchData } = useBusinessInfo({ isAdmin });
    const approveBusiness = useApproveBusiness();
    const { uploadImage } = useUploadImage();
    const { updateBusiness } = useUpdateBusiness();
    const [businessAdd, setBusinessAdd] = useState(false);
    const [isReadonly, setIsReadonly] = useState<boolean[]>([]);
    const [newData, setNewData] = useState<BusinessInfo[]>([]);
    const [timestamp, setTimestamp] = useState(Date.now())

    const delay = (ms: number) => new Promise(resolve => setTimeout(resolve, ms));

    const addNewBusiness = () => {
        setBusinessAdd(true); // Change state to show the Add Business page
    };

    const handleEditBusiness = (index: number) => {
        const newIsReadonly = [...isReadonly];
        newIsReadonly[index] = false;
        setIsReadonly(newIsReadonly);
    }

    const handleCancel = (index: number) => {
        const newIsReadonly = [...isReadonly];
        newIsReadonly[index] = true;
        setIsReadonly(newIsReadonly);
    }

    const handleInputChange = (index: number, key: string, value: string) => {
        const newData = [...businessInfo];
        newData[index][key] = value;
        setNewData(newData);
    }

    const handleSaveBusiness = (index: number) => {
        updateBusiness(newData[index]);
        const newIsReadonly = [...isReadonly];
        newIsReadonly[index] = true;
        setIsReadonly(newIsReadonly);
    }

    const handleApproveBusiness = async (id: string) => {
        approveBusiness(id);
        await delay(100);
        await fetchData();
    }

    const handleImageChange = async (event, business_id: string) => {
        const file = event.target.files[0];

        await uploadImage(file, business_id)
        setTimestamp(Date.now)
        console.log("Set false")
        // await delay(1000);
        await fetchData();
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
        return <AddBusiness setBusinessAdd={setBusinessAdd} fetchData={fetchData} />
    }

    return (
        <div className="content">
            <div className="Flex" style={{ justifyContent: "space-between", alignItems: "center" }}>
                <H2 style={{ marginBottom: "0px" }}>Manage Businesses</H2>
                <Button
                    intent="primary"
                    large

                    onClick={() => addNewBusiness()}
                >
                    Add Business
                </Button>
            </div>

            <div style={{ display: 'flex', flexDirection: 'column', gap: '20px' }}>
                {sortedBusinessInfo != null && sortedBusinessInfo.map((entry, index) => (
                    <Card interactive={true} >
                        <div className="Flex" style={{ justifyContent: "space-between" }}>
                            <H3>Business Information</H3>
                            <Tag round intent={entry.status === "active" ? Intent.SUCCESS : Intent.WARNING}>{entry.status.toLocaleUpperCase()}</Tag>
                        </div>
                        {entry.logo_url ? <img style={{ maxHeight: "40px" }} src={`${entry.logo_url}?timestamp=${timestamp.toString()}`} alt=""></img> : <Icon icon="office" />}

                        <FormGroup inline label="Name"
                            labelFor="name" >
                            <InputGroup asyncControl id="name" onChange={(e) => handleInputChange(index, "name", e.target.value)} value={entry.name} readOnly={isReadonly[index]} />
                        </FormGroup>
                        <FormGroup inline label="Website"
                            labelFor="website" >
                            <InputGroup asyncControl id="website" onChange={(e) => handleInputChange(index, "website", e.target.value)} value={entry.website} readOnly={isReadonly[index]} />
                        </FormGroup>
                        <FormGroup inline label="Image"
                            labelFor="image" >
                            <FileInput text="Choose an image..." inputProps={{
                                accept: ".png,.jpg,.jpeg", // Restrict to image files (png, jpg, jpeg)
                                onChange: (e) => handleImageChange(e, entry.id),
                            }} />
                        </FormGroup>
                        <FormGroup label="Description"
                            labelFor="desc" >
                            <TextArea autoResize asyncControl id="desc" onChange={(e) => handleInputChange(index, "desc", e.target.value)} value={entry.desc} readOnly={isReadonly[index]} fill />
                        </FormGroup>
                        {(isAdmin && entry.status === "pending") && <Button intent="success" style={{ marginRight: "10px" }} onClick={() => handleApproveBusiness(entry.id)}>Approve business</Button>}
                        {isReadonly[index] ? <Button onClick={() => handleEditBusiness(index)}>Manage business</Button> : <Button onClick={() => handleSaveBusiness(index)}>Save changes</Button>}
                        {!isReadonly[index] && <Button style={{ marginLeft: "10px" }} onClick={() => handleCancel(index)}>Cancel</Button>}
                    </Card>
                ))}
            </div>
        </div>
    );
}
