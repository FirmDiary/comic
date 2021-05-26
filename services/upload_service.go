package services

import (
    "comic/common"
    "comic/datamodels"
    "comic/repositories"
    "fmt"
    "io"
    "log"
    "mime/multipart"
    "os"
)

var dir, _ = os.Getwd()

type TransferConst struct {
    Shell     string //shell执行命令
    ShellPath string //shell执行目录
    TTL       int64  //文件存储时间 秒
}

const (
    In      = "/common/U-2-Net/upload/in/"  //文件输入目录
    Out     = "/common/U-2-Net/upload/out/" //文件输出目录
    ImgType = ".png"
)

func loadConst(transferType int) (consts *TransferConst) {
    if transferType == datamodels.TypeFace {
        consts = &TransferConst{
            Shell:     "./u2net_portrait.py",
            ShellPath: "/common/U-2-Net/",
            TTL:       60,
        }
    } else if transferType == datamodels.TypeAll {
        consts = &TransferConst{
            Shell:     "./u2net_portrait_all.py",
            ShellPath: "/common/U-2-Net/",
            TTL:       60,
        }
    }
    return
}

type IUploadService interface {
    Transfer(file multipart.File, userId int64, transferType int) (path string, err error)
}

type UploadService struct {
    repository repositories.IUploadRepository
}

func NewUploadService() IUploadService {
    return &UploadService{repository: repositories.NewUploadRepository()}
}

func (u UploadService) Transfer(file multipart.File, userId int64, transferType int) (path string, err error) {
    name := common.GetRandomString(12)

    consts := loadConst(transferType)

    fmt.Println("当前路径：", dir)
    out, err := os.OpenFile(dir+In+name+ImgType, os.O_WRONLY|os.O_CREATE, 06666)
    defer out.Close()

    if err != nil {
        return
    }

    io.Copy(out, file)

    res, err := common.CmdAndChangeDir(dir+consts.ShellPath, "python", []string{
        consts.Shell,
        name + ImgType,
    })
    if err != nil {
        return
    }

    fmt.Print(res)
    path = dir + Out + name + ImgType

    //添加数据库记录
    u.repository.Create(&datamodels.Upload{
        File:   name,
        UserId: userId,
        Type: transferType,
    })

    DelImg(name, consts.TTL)

    return
}

//删除识别的文件
func DelUploadImg(name string) error {
    fileIn := dir + In + name + ImgType
    fileOut := dir + Out + name + ImgType

    var err error
    err = os.Remove(fileIn)
    if err != nil {
        log.Fatalln(err)
        return err
    }
    err = os.Remove(fileOut)
    if err != nil {
        log.Fatalln(err)
        return err
    }
    return nil
}
